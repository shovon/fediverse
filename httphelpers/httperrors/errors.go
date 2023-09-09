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
		status: http.StatusBadRequest,
		err:    "Bad Request",
	}
}

func Unauthorized() HTTPError {
	return dummy{
		status: http.StatusUnauthorized,
		err:    "Unauthorized",
	}
}

func Forbidden() HTTPError {
	return dummy{
		status: http.StatusForbidden,
		err:    "Forbidden",
	}
}

func NotFound() HTTPError {
	return dummy{
		status: http.StatusNotFound,
		err:    "Not Found",
	}
}

func NotAcceptable() HTTPError {
	return dummy{
		status: http.StatusNotAcceptable,
		err:    "Not Acceptable",
	}
}

func ProxyAuthenticationRequired() HTTPError {
	return dummy{
		status: http.StatusProxyAuthRequired,
		err:    "Proxy Authentication Required",
	}
}

func RequestTimeout() HTTPError {
	return dummy{
		status: http.StatusRequestTimeout,
		err:    "Request Timeout",
	}
}

func Conflict() HTTPError {
	return dummy{
		status: http.StatusConflict,
		err:    "Conflict",
	}
}

func Gone() HTTPError {
	return dummy{
		status: http.StatusGone,
		err:    "Gone",
	}
}

func LengthRequired() HTTPError {
	return dummy{
		status: http.StatusLengthRequired,
		err:    "Length Required",
	}
}

func PreconditionFailed() HTTPError {
	return dummy{
		status: http.StatusPreconditionFailed,
		err:    "Precondition Failed",
	}
}

func RequestEntityTooLarge() HTTPError {
	return dummy{
		status: http.StatusRequestEntityTooLarge,
		err:    "Request Entity Too Large",
	}
}

func RequestURITooLong() HTTPError {
	return dummy{
		status: http.StatusRequestURITooLong,
		err:    "Request URI Too Long",
	}
}

func UnsupportedMediaType() HTTPError {
	return dummy{
		status: http.StatusUnsupportedMediaType,
		err:    "Unsupported Media Type",
	}
}

func RequestedRangeNotSatisfiable() HTTPError {
	return dummy{
		status: http.StatusRequestedRangeNotSatisfiable,
		err:    "Requested Range Not Satisfiable",
	}
}

func ExpectationFailed() HTTPError {
	return dummy{
		status: http.StatusExpectationFailed,
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
		status: http.StatusMisdirectedRequest,
		err:    "Misdirected Request",
	}
}

func UnprocessableEntity() HTTPError {
	return dummy{
		status: http.StatusUnprocessableEntity,
		err:    "Unprocessable Entity",
	}
}

func Locked() HTTPError {
	return dummy{
		status: http.StatusLocked,
		err:    "Locked",
	}
}

func FailedDependency() HTTPError {
	return dummy{
		status: http.StatusFailedDependency,
		err:    "Failed Dependency",
	}
}

func TooEarly() HTTPError {
	return dummy{
		status: http.StatusTooEarly,
		err:    "Too Early",
	}
}

func UpgradeRequired() HTTPError {
	return dummy{
		status: http.StatusUpgradeRequired,
		err:    "Upgrade Required",
	}
}

func PreconditionRequired() HTTPError {
	return dummy{
		status: http.StatusPreconditionRequired,
		err:    "Precondition Required",
	}
}

func TooManyRequests() HTTPError {
	return dummy{
		status: http.StatusTooManyRequests,
		err:    "Too Many Requests",
	}
}

func RequestHeaderFieldsTooLarge() HTTPError {
	return dummy{
		status: http.StatusRequestHeaderFieldsTooLarge,
		err:    "Request Header Fields Too Large",
	}
}

func UnavailableForLegalReasons() HTTPError {
	return dummy{
		status: http.StatusUnavailableForLegalReasons,
		err:    "Unavailable For Legal Reasons",
	}
}

func InternalServerError() HTTPError {
	return dummy{
		status: http.StatusInternalServerError,
		err:    "Internal Server Error",
	}
}

func NotImplemented() HTTPError {
	return dummy{
		status: http.StatusNotImplemented,
		err:    "Not Implemented",
	}
}

func BadGateway() HTTPError {
	return dummy{
		status: http.StatusBadGateway,
		err:    "Bad Gateway",
	}
}

func ServiceUnavailable() HTTPError {
	return dummy{
		status: http.StatusServiceUnavailable,
		err:    "Service Unavailable",
	}
}

func GatewayTimeout() HTTPError {
	return dummy{
		status: http.StatusGatewayTimeout,
		err:    "Gateway Timeout",
	}
}

func HTTPVersionNotSupported() HTTPError {
	return dummy{
		status: http.StatusHTTPVersionNotSupported,
		err:    "HTTP Version Not Supported",
	}
}

func VariantAlsoNegotiates() HTTPError {
	return dummy{
		status: http.StatusVariantAlsoNegotiates,
		err:    "Variant Also Negotiates",
	}
}

func InsufficientStorage() HTTPError {
	return dummy{
		status: http.StatusInsufficientStorage,
		err:    "Insufficient Storage",
	}
}

func LoopDetected() HTTPError {
	return dummy{
		status: http.StatusLoopDetected,
		err:    "Loop Detected",
	}
}

func NotExtended() HTTPError {
	return dummy{
		status: http.StatusNotExtended,
		err:    "Not Extended",
	}
}

func NetworkAuthenticationRequired() HTTPError {
	return dummy{
		status: http.StatusNetworkAuthenticationRequired,
		err:    "Network Authentication Required",
	}
}
