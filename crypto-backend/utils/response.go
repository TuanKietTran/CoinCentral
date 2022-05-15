package utils

import (
	"log"
	"net/http"
)

func CloseResponseBody(resp *http.Response) {
	if err := resp.Body.Close(); err != nil {
		log.Printf("Can't close response's body, just skipping, err: %v\n", err)
	}
}
