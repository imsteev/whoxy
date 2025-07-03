package main

import (
	"bytes"
	"net/http"
)

// ForwardPostRequest forwards a POST request given an existing request. Headers will be copied.
func ForwardPostRequest(r *http.Request, dest string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", dest, bytes.NewReader(body)) // do i need to close the body?
	if err != nil {
		return nil, err
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return req, nil
}
