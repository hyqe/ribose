package fit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hyqe/ribose/internal/fit/status"
)

type RPC struct {
	name    string // name of rpc struct
	methods map[string]*Method
	ptr     reflect.Value
	*validator.Validate
}

// NewRPC builds an RPC with from an instance of a type and
// its methods.
//
// The exported methods must have the following signature.
// Where INPUT is a pointer to your input type, and OUTPUT
// is value your method sends back to the client.
//
//	fn(ctx context.Context, in *INPUT) (*OUTPUT, status.Status)
//
// INPUT is validated before being passed to it method. see
// https://pkg.go.dev/github.com/go-playground/validator/v10
//
//	type Foo struct {
//		Email string `validate:"email"`
//	}
//
// Docs endpoints are created for each method.
//
//	GET /<type>/help // gets list of methods
//	GET /<type>/<method>/help // gets INPUT/OUTPUT
func NewRPC(ptr any) *RPC {
	reflectVal := reflect.ValueOf(ptr)
	return &RPC{
		name:     getTypeName(ptr),
		methods:  parseMethods(reflectVal),
		ptr:      reflectVal,
		Validate: validator.New(),
	}
}

func parseMethods(v reflect.Value) map[string]*Method {
	methods := make(map[string]*Method)
	methodNames := ListExportedMethodNames(v.Type())
	for _, methodName := range methodNames {
		if method, ok := newMethod(v, methodName); ok {
			methods[methodName] = &method
		}
	}
	return methods
}

func (s *RPC) docsJSON() any {
	methods := make([]string, 0, len(s.methods))
	for _, m := range s.methods {
		methods = append(methods, m.name)
	}
	return map[string]any{
		"methods": methods,
	}
}

func (s *RPC) MountFiberApp(app *fiber.App) fiber.Router {
	sub := fiber.New()

	sub.Get("/help", func(c *fiber.Ctx) error {
		return c.JSON(s.docsJSON())
	})

	for _, m := range s.methods {
		func(method *Method) {
			subPath, _ := url.JoinPath("/", method.name)
			subHelp, _ := url.JoinPath("/", method.name, "help")
			sub.Get(path.Join(subHelp), func(c *fiber.Ctx) error {
				return c.JSON(method.docsJSON())
			})
			sub.Post(subPath, func(c *fiber.Ctx) error {
				in := method.NewIn().Interface()
				err := c.BodyParser(in)
				if err != nil {
					return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("failed to decode body: %v", err))
				}
				err = s.Validate.Struct(in)
				if err != nil {
					return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("validation failed: %v", err))
				}
				out, status := method.Invoke(c.Context(), in)
				if status.Code >= 300 {
					return fiber.NewError(int(status.Code), status.Message)
				}
				err = c.JSON(out)
				if err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to encode response: %v", err))
				}
				return c.SendStatus(int(status.Code))
			})
		}(m)
	}

	return app.Mount("/"+s.name, sub)
}

func (s *RPC) NewNetHttpHandler() http.HandlerFunc {
	mux := http.NewServeMux()
	for _, m := range s.methods {
		func(method *Method) {
			fullPath, _ := url.JoinPath("/", s.name, method.name)
			mux.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
				in := method.NewIn().Interface()
				err := json.NewDecoder(r.Body).Decode(in)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to decode body as json: %v", err), http.StatusBadRequest)
					return
				}
				err = s.Validate.Struct(in)
				if err != nil {
					http.Error(w, fmt.Sprintf("validation failed: %v", err), http.StatusBadRequest)
					return
				}
				out, status := method.Invoke(r.Context(), in)
				if status.Code >= 300 {
					http.Error(w, status.Message, int(status.Code))
					return
				}
				w.WriteHeader(int(status.Code))
				json.NewEncoder(w).Encode(out)
			})
		}(m)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}
}

type Method struct {
	In      string
	inType  reflect.Type
	Out     string
	outType reflect.Type

	svc  reflect.Value
	name string
	fn   reflect.Value
}

// NewMethod is looking for a method with this signature:
//
//	fn[I, O any](ctx context.Context, in I) (O, status.Status)
func newMethod(svc reflect.Value, methodName string) (method Method, ok bool) {
	parentReflectedType := svc.Type()
	reflectedMethod, _ := parentReflectedType.MethodByName(methodName)
	// (ctx context.Context, in I)
	NumIn := reflectedMethod.Type.NumIn()
	if NumIn != 3 {
		return
	}

	// first arg is a context
	reflectContext := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !reflectedMethod.Type.In(1).Implements(reflectContext) {
		return
	}

	// second arg is a *struct
	if reflect.Pointer != reflectedMethod.Type.In(2).Kind() {
		return
	}
	if reflect.Struct != reflectedMethod.Type.In(2).Elem().Kind() {
		return
	}

	// (O, status.Status)
	NumOut := reflectedMethod.Type.NumOut()
	if NumOut != 2 {
		return
	}

	// first return is a *struct
	if reflect.Pointer != reflectedMethod.Type.Out(0).Kind() {
		return
	}
	if reflect.Struct != reflectedMethod.Type.Out(0).Elem().Kind() {
		return
	}

	// second return is a status.Status
	if reflect.TypeOf(status.Status{}) != reflectedMethod.Type.Out(1) {
		return
	}

	return Method{
		In:      reflectedMethod.Type.In(2).Name(),
		inType:  reflectedMethod.Type.In(2),
		Out:     reflectedMethod.Type.Out(0).Name(),
		outType: reflectedMethod.Type.Out(0),

		svc:  svc,
		name: reflectedMethod.Name,
		fn:   reflectedMethod.Func,
	}, true
}

func (m *Method) docsJSON() any {
	return map[string]any{
		"request": Property{
			Type:       "object",
			Properties: buildStructDocs(m.inType.Elem()),
		},
		"response": Property{
			Type:       "object",
			Properties: buildStructDocs(m.outType.Elem()),
		},
	}
}

// buildStructDocs support maps/slices fields
func buildStructDocs(v reflect.Type) map[string]Property {
	out := make(map[string]Property)
	NumField := v.NumField()
	for f := 0; f < NumField; f++ {
		field := v.Field(f)

		if !field.IsExported() {
			continue
		}

		var fieldName string
		jsonTag, ok := field.Tag.Lookup("json")
		if ok {
			jsonTagParts := strings.Split(jsonTag, ",")
			fieldName = jsonTagParts[0]
		}

		var example string
		exampleTag, ok := field.Tag.Lookup("example")
		if ok {
			example = exampleTag
		}

		var format string
		formatTag, ok := field.Tag.Lookup("format")
		if ok {
			format = formatTag
		}

		//var validate string
		validateTag, ok := field.Tag.Lookup("validate")
		if ok {
			format = validateTag
		}

		var fieldTypeName string
		switch field.Type {
		case reflect.TypeOf(uuid.UUID{}):
			out[fieldName] = Property{
				Type:    "string",
				Format:  "uuid",
				Example: uuid.NewString(),
			}
		case reflect.TypeOf(time.Time{}):
			out[fieldName] = Property{
				Type:    "string",
				Format:  "rfc3339",
				Example: time.Now().UTC().Format(time.RFC3339),
			}
		default:
			switch field.Type.Kind() {
			case reflect.Struct:
				out[fieldName] = Property{
					Type:       "object",
					Properties: buildStructDocs(field.Type),
					Example:    example,
					//Validate:   validate,
					Format: format,
				}
			case reflect.Map:
				out[fieldName] = Property{
					Type:    "object",
					Example: example,
					//Validate: validate,
					Format: format,
				}
			case reflect.Slice:
				out[fieldName] = Property{
					Type:    "array",
					Example: example,
					//Validate: validate,
					Format: format,
				}
			default:
				fieldTypeName = field.Type.Kind().String()
				out[fieldName] = Property{
					Type:    fieldTypeName,
					Example: example,
					//Validate: validate,
					Format: format,
				}
			}
		}

	}
	return out
}

type Property struct {
	Type       string              `json:"type,omitempty"`
	Format     string              `json:"format,omitempty"`
	Validate   string              `json:"validate,omitempty"`
	Example    string              `json:"example,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// NewIn mints a new inType.
//
//	in := m.NewIn()
//	json.Unmarshal(data, in.Interface())
func (m *Method) NewIn() reflect.Value {
	if m.inType.Kind() == reflect.Pointer {
		return reflect.New(m.inType.Elem())
	}
	return reflect.New(m.inType)
}

func (m *Method) Invoke(ctx context.Context, in any) (any, status.Status) {
	resp := m.fn.Call([]reflect.Value{m.svc, reflect.ValueOf(ctx), reflect.ValueOf(in)})

	out := resp[0].Interface()
	status := resp[1].Interface().(status.Status)
	return out, status
}

func getTypeName(v any) string {
	return reflect.Indirect(reflect.ValueOf(v)).Type().Name()
}

func ListExportedMethodNames(v reflect.Type) []string {
	methodNames := make([]string, 0)

	numberOfMethods := v.NumMethod()
	for m := 0; m < numberOfMethods; m++ {
		reflectedMethod := v.Method(m)

		// Method must be exported.
		if !reflectedMethod.IsExported() {
			continue
		}

		methodNames = append(methodNames, reflectedMethod.Name)
	}

	return methodNames
}
