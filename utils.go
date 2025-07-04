package main

import (
	"context"
	"net/http"
	"net/url"
)

func ForwardPostRequest(r *http.Request, destination string) (*http.Response, error) {
	destinationUrl, err := url.Parse(destination)
	if err != nil {
		return nil, err
	}
	clonedReq := r.Clone(context.Background())
	clonedReq.URL = destinationUrl
	return http.DefaultClient.Do(clonedReq)
}
