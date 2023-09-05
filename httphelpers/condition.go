package httphelpers

import (
	"context"
	"net/http"
	"net/url"
)

type BarebonesRequest struct {
	originalRequest  *http.Request
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

func (r BarebonesRequest) Context() context.Context {
	return r.originalRequest.Context()
}

func CopyRequest(req *http.Request) (BarebonesRequest, error) {
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return BarebonesRequest{}, err
	}
	return BarebonesRequest{
		originalRequest:  req,
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
				reqCopy, err := CopyRequest(r)
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