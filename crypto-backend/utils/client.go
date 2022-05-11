package utils

import (
	"net/http"
	"time"
)

// HttpClient is used to share HTTP Client between our packages
var HttpClient = &http.Client{Timeout: 50 * time.Second}

func init() {
	newTransport := http.DefaultTransport.(*http.Transport).Clone()
	newTransport.MaxIdleConns = 100
	newTransport.MaxConnsPerHost = 100
	newTransport.MaxIdleConnsPerHost = 100

	HttpClient.Transport = newTransport
}
