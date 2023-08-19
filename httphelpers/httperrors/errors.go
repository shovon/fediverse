package httperrors

import "net/http"

type useless struct{}

type HTTPError interface {
	useless() useless
	Status() int
	Error() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type dummy struct {
	status int
	err    string
}

var _ HTTPError = dummy{}
var _ error = dummy{}
var _ http.Handler = dummy{}

func (d dummy) useless() useless {
	return useless{}
}

func (d dummy) Status() int {
	return d.status
}

func (d dummy) Error() string {
	return d.err
}

func (d dummy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(d.status)
	w.Write([]byte(d.err))
}

func BadRequest() HTTPError {
	return dummy{
		status: 400,
		err:    "Bad Request",
	}
}

func Unauthorized() HTTPError {
	return dummy{
		status: 401,
		err:    "Unauthorized",
	}
}

func Forbidden() HTTPError {
	return dummy{
		status: 403,
		err:    "Forbidden",
	}
}

func NotFound() HTTPError {
	return dummy{
		status: 404,
		err:    "Not Found",
	}
}

func NotAcceptable() HTTPError {
	return dummy{
		status: 406,
		err:    "Not Acceptable",
	}
}

func ProxyAuthenticationRequired() HTTPError {
	return dummy{
		status: 407,
		err:    "Proxy Authentication Required",
	}
}

func RequestTimeout() HTTPError {
	return dummy{
		status: 408,
		err:    "Request Timeout",
	}
}

func Conflict() HTTPError {
	return dummy{
		status: 409,
		err:    "Conflict",
	}
}

func Gone() HTTPError {
	return dummy{
		status: 410,
		err:    "Gone",
	}
}

func LengthRequired() HTTPError {
	return dummy{
		status: 411,
		err:    "Length Required",
	}
}

func PreconditionFailed() HTTPError {
	return dummy{
		status: 412,
		err:    "Precondition Failed",
	}
}

func RequestEntityTooLarge() HTTPError {
	return dummy{
		status: 413,
		err:    "Request Entity Too Large",
	}
}

func RequestURITooLong() HTTPError {
	return dummy{
		status: 414,
		err:    "Request URI Too Long",
	}
}

func UnsupportedMediaType() HTTPError {
	return dummy{
		status: 415,
		err:    "Unsupported Media Type",
	}
}

func RequestedRangeNotSatisfiable() HTTPError {
	return dummy{
		status: 416,
		err:    "Requested Range Not Satisfiable",
	}
}

func ExpectationFailed() HTTPError {
	return dummy{
		status: 417,
		err:    "Expectation Failed",
	}
}

func ImATeapot() HTTPError {
	return dummy{
		status: 418,
		err:    "I'm a teapot",
	}
}

func MisdirectedRequest() HTTPError {
	return dummy{
		status: 421,
		err:    "Misdirected Request",
	}
}

func UnprocessableEntity() HTTPError {
	return dummy{
		status: 422,
		err:    "Unprocessable Entity",
	}
}

func Locked() HTTPError {
	return dummy{
		status: 423,
		err:    "Locked",
	}
}

func FailedDependency() HTTPError {
	return dummy{
		status: 424,
		err:    "Failed Dependency",
	}
}

func TooEarly() HTTPError {
	return dummy{
		status: 425,
		err:    "Too Early",
	}
}

func UpgradeRequired() HTTPError {
	return dummy{
		status: 426,
		err:    "Upgrade Required",
	}
}

func PreconditionRequired() HTTPError {
	return dummy{
		status: 428,
		err:    "Precondition Required",
	}
}

func TooManyRequests() HTTPError {
	return dummy{
		status: 429,
		err:    "Too Many Requests",
	}
}

func RequestHeaderFieldsTooLarge() HTTPError {
	return dummy{
		status: 431,
		err:    "Request Header Fields Too Large",
	}
}

func UnavailableForLegalReasons() HTTPError {
	return dummy{
		status: 451,
		err:    "Unavailable For Legal Reasons",
	}
}

func InternalServerError() HTTPError {
	return dummy{
		status: 500,
		err:    "Internal Server Error",
	}
}

func NotImplemented() HTTPError {
	return dummy{
		status: 501,
		err:    "Not Implemented",
	}
}

func BadGateway() HTTPError {
	return dummy{
		status: 502,
		err:    "Bad Gateway",
	}
}

func ServiceUnavailable() HTTPError {
	return dummy{
		status: 503,
		err:    "Service Unavailable",
	}
}

func GatewayTimeout() HTTPError {
	return dummy{
		status: 504,
		err:    "Gateway Timeout",
	}
}

func HTTPVersionNotSupported() HTTPError {
	return dummy{
		status: 505,
		err:    "HTTP Version Not Supported",
	}
}

func VariantAlsoNegotiates() HTTPError {
	return dummy{
		status: 506,
		err:    "Variant Also Negotiates",
	}
}

func InsufficientStorage() HTTPError {
	return dummy{
		status: 507,
		err:    "Insufficient Storage",
	}
}

func LoopDetected() HTTPError {
	return dummy{
		status: 508,
		err:    "Loop Detected",
	}
}

func NotExtended() HTTPError {
	return dummy{
		status: 510,
		err:    "Not Extended",
	}
}

func NetworkAuthenticationRequired() HTTPError {
	return dummy{
		status: 511,
		err:    "Network Authentication Required",
	}
}
