package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/imsteev/whoxy/integrations/clerk"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("whoxy"))
	})

	http.HandleFunc("/clerk", func(w http.ResponseWriter, r *http.Request) {
		clerkIntegration := &clerk.ClerkIntegration{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		dest, err := clerkIntegration.GetDestination(body)
		if err != nil {
			log.Fatal(err)
		}

		// Create a new request with headers
		req, err := http.NewRequest("POST", dest, bytes.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}

		// Copy headers from incoming request
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
		req.Header.Add("X-Whoxy-Integration", "clerk")

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		w.Write([]byte("OK"))
	})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
