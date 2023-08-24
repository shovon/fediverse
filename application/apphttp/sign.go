package apphttp

import "net/http"

func SendSigned(req *http.Request, key any, body string) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}
