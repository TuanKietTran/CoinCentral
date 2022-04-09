package utils

import "net/http"

// HttpClient is used to share HTTP Client between our packages
var HttpClient = &http.Client{}
