package fit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hyqe/ribose/internal/fit/status"
)

type RPC struct {
	name    string // name of rpc struct
	methods map[string]*Method
	ptr     reflect.Value
	*validator.Validate
}

func NewRPC(ptr any) *RPC {
	reflectVal := reflect.ValueOf(ptr)
	methods := make(map[string]*Method)

	methodNames := ListExportedMethodNames(ptr)

	for _, methodName := range methodNames {
		if method, ok := newMethod(reflectVal, methodName); ok {
			methods[methodName] = &method
		}
	}

	return &RPC{
		name:     getTypeName(ptr),
		methods:  methods,
		ptr:      reflectVal,
		Validate: validator.New(),
	}
}

func (s *RPC) MountFiberApp(app *fiber.App) fiber.Router {
	sub := fiber.New()

	for _, m := range s.methods {
		func(method *Method) {
			sub.Post("/"+method.name, func(c *fiber.Ctx) error {
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
			mux.HandleFunc(path.Join("/", s.name, method.name), func(w http.ResponseWriter, r *http.Request) {
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
	//
	// TODO: find the correct way to check the type of an context.Context.
	// since context.Context is an interface, its not possible to just mint a new
	// instance of it, and check its interface.

	firstArgTypName := reflectedMethod.Type.In(1).Name()
	firstArgTypPath := reflectedMethod.Type.In(1).PkgPath()
	firstArgFullTypName := firstArgTypPath + "." + firstArgTypName
	if firstArgFullTypName != "context.Context" {
		return
	}

	// second arg is a struct
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

	// second arg is a struct
	if reflect.Pointer != reflectedMethod.Type.Out(0).Kind() {
		return
	}
	if reflect.Struct != reflectedMethod.Type.Out(0).Elem().Kind() {
		return
	}

	// second arg is a struct
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

func ListExportedMethodNames(v any) []string {
	methodNames := make([]string, 0)

	reflectedType := reflect.TypeOf(v)

	numberOfMethods := reflectedType.NumMethod()
	for m := 0; m < numberOfMethods; m++ {
		reflectedMethod := reflectedType.Method(m)

		// Method must be exported.
		if !reflectedMethod.IsExported() {
			continue
		}

		methodNames = append(methodNames, reflectedMethod.Name)
	}

	return methodNames
}
