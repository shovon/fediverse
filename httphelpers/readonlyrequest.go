package httphelpers

import (
	"context"
	"net/http"
	"net/url"
)

// ReadOnlyRequest is a stripped-down version of http.Request, so that we don't
// end up writing silly things to it.
type ReadOnlyRequest struct {
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

type SomeSubset struct {
	Method string
	URL    *url.URL
}

func GetSomeSomeSubset(r ReadOnlyRequest) SomeSubset {
	return SomeSubset{
		Method: r.Method,
		URL:    r.URL,
	}
}

// Context returns the context of the original request.
func (r ReadOnlyRequest) Context() context.Context {
	return r.originalRequest.Context()
}

// ToReadOnlyRequest converts an http.Request to a ReadOnlyRequest.
func ToReadOnlyRequest(req *http.Request) (ReadOnlyRequest, error) {
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return ReadOnlyRequest{}, err
	}
	return ReadOnlyRequest{
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
