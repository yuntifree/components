package httputil

import (
	"io/ioutil"
	"net/http"
	"strings"
)

//Request return response body of http request
func Request(url, reqbody string) (string, error) {
	return RequestWithHeaders(url, reqbody, map[string]string{})
}

//RequestWithHeaders return response body of http request with headers
func RequestWithHeaders(url, reqbody string, headers map[string]string) (string, error) {
	client := &http.Client{}
	method := "GET"
	if len(reqbody) > 0 {
		method = "POST"
	}
	req, err := http.NewRequest(method, url, strings.NewReader(reqbody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	rspbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(rspbody), nil
}
