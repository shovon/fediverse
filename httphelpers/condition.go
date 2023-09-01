package httphelpers

import (
	"net/http"
	"net/url"
)

type BarebonesRequest struct {
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	ContentLength    int64
	TransferEncoding []string
	Host             string
	RemoteAddr       string
	RequestURI       string
}

func copyRequest(req *http.Request) (BarebonesRequest, error) {
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return BarebonesRequest{}, err
	}
	return BarebonesRequest{
		Method:           req.Method,
		URL:              u,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Header:           req.Header,
		ContentLength:    req.ContentLength,
		TransferEncoding: req.TransferEncoding,
		Host:             req.Host,
		RemoteAddr:       req.RemoteAddr,
		RequestURI:       req.RequestURI,
	}, nil
}

func Condition(predicate func(r BarebonesRequest) bool) Processor {
	return ProcessorFunc(func(middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqCopy, err := copyRequest(r)
				if err != nil {
					middleware(next).ServeHTTP(w, r)
					return
				}
				if predicate(reqCopy) {
					middleware(next).ServeHTTP(w, r)
					return
				}
				next.ServeHTTP(w, r)
			})
		}
	})
}
