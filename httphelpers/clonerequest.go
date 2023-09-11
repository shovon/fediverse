package httphelpers

import (
	"net/http"
	"net/url"
)

func CloneHRequest(req *http.Request) (*http.Request, error) {
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, err
	}
	return &http.Request{
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
