package jrdhttperrors

import "net/http"

type useless struct{}

type JRDHttpError interface {
	Useless() useless
	Status() int
	Error() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type dummy struct {
	status int
	err    string
}

var _ JRDHttpError = dummy{}
var _ error = dummy{}
var _ http.Handler = dummy{}

func (d dummy) Useless() useless {
	return useless{}
}

func (d dummy) Status() int {
	return d.status
}

func (d dummy) Error() string {
	return d.err
}

func (d dummy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(d.status)
	w.Write([]byte(d.err))
}

func BadRequest() JRDHttpError {
	return dummy{
		status: 400,
		err:    "Bad Request",
	}
}

func Unauthorized() JRDHttpError {
	return dummy{
		status: 401,
		err:    "Unauthorized",
	}
}

func Forbidden() JRDHttpError {
	return dummy{
		status: 403,
		err:    "Forbidden",
	}
}

func NotFound() JRDHttpError {
	return dummy{
		status: 404,
		err:    "Not Found",
	}
}

func NotAcceptable() JRDHttpError {
	return dummy{
		status: 406,
		err:    "Not Acceptable",
	}
}

func ProxyAuthenticationRequired() JRDHttpError {
	return dummy{
		status: 407,
		err:    "Proxy Authentication Required",
	}
}

func RequestTimeout() JRDHttpError {
	return dummy{
		status: 408,
		err:    "Request Timeout",
	}
}

func Conflict() JRDHttpError {
	return dummy{
		status: 409,
		err:    "Conflict",
	}
}

func Gone() JRDHttpError {
	return dummy{
		status: 410,
		err:    "Gone",
	}
}

func LengthRequired() JRDHttpError {
	return dummy{
		status: 411,
		err:    "Length Required",
	}
}

func PreconditionFailed() JRDHttpError {
	return dummy{
		status: 412,
		err:    "Precondition Failed",
	}
}

func RequestEntityTooLarge() JRDHttpError {
	return dummy{
		status: 413,
		err:    "Request Entity Too Large",
	}
}

func RequestURITooLong() JRDHttpError {
	return dummy{
		status: 414,
		err:    "Request URI Too Long",
	}
}

func UnsupportedMediaType() JRDHttpError {
	return dummy{
		status: 415,
		err:    "Unsupported Media Type",
	}
}

func RequestedRangeNotSatisfiable() JRDHttpError {
	return dummy{
		status: 416,
		err:    "Requested Range Not Satisfiable",
	}
}

func ExpectationFailed() JRDHttpError {
	return dummy{
		status: 417,
		err:    "Expectation Failed",
	}
}

func ImATeapot() JRDHttpError {
	return dummy{
		status: 418,
		err:    "I'm a teapot",
	}
}

func MisdirectedRequest() JRDHttpError {
	return dummy{
		status: 421,
		err:    "Misdirected Request",
	}
}

func UnprocessableEntity() JRDHttpError {
	return dummy{
		status: 422,
		err:    "Unprocessable Entity",
	}
}

func Locked() JRDHttpError {
	return dummy{
		status: 423,
		err:    "Locked",
	}
}

func FailedDependency() JRDHttpError {
	return dummy{
		status: 424,
		err:    "Failed Dependency",
	}
}

func TooEarly() JRDHttpError {
	return dummy{
		status: 425,
		err:    "Too Early",
	}
}

func UpgradeRequired() JRDHttpError {
	return dummy{
		status: 426,
		err:    "Upgrade Required",
	}
}

func PreconditionRequired() JRDHttpError {
	return dummy{
		status: 428,
		err:    "Precondition Required",
	}
}

func TooManyRequests() JRDHttpError {
	return dummy{
		status: 429,
		err:    "Too Many Requests",
	}
}

func RequestHeaderFieldsTooLarge() JRDHttpError {
	return dummy{
		status: 431,
		err:    "Request Header Fields Too Large",
	}
}

func UnavailableForLegalReasons() JRDHttpError {
	return dummy{
		status: 451,
		err:    "Unavailable For Legal Reasons",
	}
}

func InternalServerError() JRDHttpError {
	return dummy{
		status: 500,
		err:    "Internal Server Error",
	}
}

func NotImplemented() JRDHttpError {
	return dummy{
		status: 501,
		err:    "Not Implemented",
	}
}

func BadGateway() JRDHttpError {
	return dummy{
		status: 502,
		err:    "Bad Gateway",
	}
}

func ServiceUnavailable() JRDHttpError {
	return dummy{
		status: 503,
		err:    "Service Unavailable",
	}
}

func GatewayTimeout() JRDHttpError {
	return dummy{
		status: 504,
		err:    "Gateway Timeout",
	}
}

func HTTPVersionNotSupported() JRDHttpError {
	return dummy{
		status: 505,
		err:    "HTTP Version Not Supported",
	}
}

func VariantAlsoNegotiates() JRDHttpError {
	return dummy{
		status: 506,
		err:    "Variant Also Negotiates",
	}
}

func InsufficientStorage() JRDHttpError {
	return dummy{
		status: 507,
		err:    "Insufficient Storage",
	}
}

func LoopDetected() JRDHttpError {
	return dummy{
		status: 508,
		err:    "Loop Detected",
	}
}

func NotExtended() JRDHttpError {
	return dummy{
		status: 510,
		err:    "Not Extended",
	}
}

func NetworkAuthenticationRequired() JRDHttpError {
	return dummy{
		status: 511,
		err:    "Network Authentication Required",
	}
}
