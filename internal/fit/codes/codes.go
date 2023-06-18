package codes

import "github.com/lib/pq"

type Code int

const (
	Continue                      Code = 100 // RFC 9110, 15.2.1
	SwitchingProtocols            Code = 101 // RFC 9110, 15.2.2
	Processing                    Code = 102 // RFC 2518, 10.1
	EarlyHints                    Code = 103 // RFC 8297
	OK                            Code = 200 // RFC 9110, 15.3.1
	Created                       Code = 201 // RFC 9110, 15.3.2
	Accepted                      Code = 202 // RFC 9110, 15.3.3
	NonAuthoritativeInfo          Code = 203 // RFC 9110, 15.3.4
	NoContent                     Code = 204 // RFC 9110, 15.3.5
	ResetContent                  Code = 205 // RFC 9110, 15.3.6
	PartialContent                Code = 206 // RFC 9110, 15.3.7
	MultiStatus                   Code = 207 // RFC 4918, 11.1
	AlreadyReported               Code = 208 // RFC 5842, 7.1
	IMUsed                        Code = 226 // RFC 3229, 10.4.1
	MultipleChoices               Code = 300 // RFC 9110, 15.4.1
	MovedPermanently              Code = 301 // RFC 9110, 15.4.2
	Found                         Code = 302 // RFC 9110, 15.4.3
	SeeOther                      Code = 303 // RFC 9110, 15.4.4
	NotModified                   Code = 304 // RFC 9110, 15.4.5
	UseProxy                      Code = 305 // RFC 9110, 15.4.6
	TemporaryRedirect             Code = 307 // RFC 9110, 15.4.8
	PermanentRedirect             Code = 308 // RFC 9110, 15.4.9
	BadRequest                    Code = 400 // RFC 9110, 15.5.1
	Unauthorized                  Code = 401 // RFC 9110, 15.5.2
	PaymentRequired               Code = 402 // RFC 9110, 15.5.3
	Forbidden                     Code = 403 // RFC 9110, 15.5.4
	NotFound                      Code = 404 // RFC 9110, 15.5.5
	MethodNotAllowed              Code = 405 // RFC 9110, 15.5.6
	NotAcceptable                 Code = 406 // RFC 9110, 15.5.7
	ProxyAuthRequired             Code = 407 // RFC 9110, 15.5.8
	RequestTimeout                Code = 408 // RFC 9110, 15.5.9
	Conflict                      Code = 409 // RFC 9110, 15.5.10
	Gone                          Code = 410 // RFC 9110, 15.5.11
	LengthRequired                Code = 411 // RFC 9110, 15.5.12
	PreconditionFailed            Code = 412 // RFC 9110, 15.5.13
	RequestEntityTooLarge         Code = 413 // RFC 9110, 15.5.14
	RequestURITooLong             Code = 414 // RFC 9110, 15.5.15
	UnsupportedMediaType          Code = 415 // RFC 9110, 15.5.16
	RequestedRangeNotSatisfiable  Code = 416 // RFC 9110, 15.5.17
	ExpectationFailed             Code = 417 // RFC 9110, 15.5.18
	Teapot                        Code = 418 // RFC 9110, 15.5.19 (Unused)
	MisdirectedRequest            Code = 421 // RFC 9110, 15.5.20
	UnprocessableEntity           Code = 422 // RFC 9110, 15.5.21
	Locked                        Code = 423 // RFC 4918, 11.3
	FailedDependency              Code = 424 // RFC 4918, 11.4
	TooEarly                      Code = 425 // RFC 8470, 5.2.
	UpgradeRequired               Code = 426 // RFC 9110, 15.5.22
	PreconditionRequired          Code = 428 // RFC 6585, 3
	TooManyRequests               Code = 429 // RFC 6585, 4
	RequestHeaderFieldsTooLarge   Code = 431 // RFC 6585, 5
	UnavailableForLegalReasons    Code = 451 // RFC 7725, 3
	InternalServerError           Code = 500 // RFC 9110, 15.6.1
	NotImplemented                Code = 501 // RFC 9110, 15.6.2
	BadGateway                    Code = 502 // RFC 9110, 15.6.3
	ServiceUnavailable            Code = 503 // RFC 9110, 15.6.4
	GatewayTimeout                Code = 504 // RFC 9110, 15.6.5
	HTTPVersionNotSupported       Code = 505 // RFC 9110, 15.6.6
	VariantAlsoNegotiates         Code = 506 // RFC 2295, 8.1
	InsufficientStorage           Code = 507 // RFC 4918, 11.5
	LoopDetected                  Code = 508 // RFC 5842, 7.2
	NotExtended                   Code = 510 // RFC 2774, 7
	NetworkAuthenticationRequired Code = 511 // RFC 6585, 6
)

// Alias
const (
	Internal      = InternalServerError
	Invalid       = BadRequest
	AlreadyExists = Conflict
)

// Pg converts an pq.ErrorCode into a Code.
// https://www.postgresql.org/docs/current/errcodes-appendix.html
// https://www.postgresql.org/docs/current/protocol-error-fields.html
func Pg(e pq.ErrorCode) Code {
	switch string(e) {
	case "23505":
		return AlreadyExists
	default:
		return Internal
	}
}
