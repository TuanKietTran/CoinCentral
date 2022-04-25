package fb

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	headerNameXSign = "X-Hub-Signature"
	signaturePrefix = "sha1="
)

// errors
var (
	errNoXSignHeader      = errors.New("there is no x-sign header")
	errInvalidXSignHeader = errors.New("invalid x-sign header")
)

// Authorize authorizes web hooks from FB
// https://developers.facebook.com/docs/graph-api/webhooks/getting-started/#validate-payloads
func Authorize(r *http.Request) error {
	signature := r.Header.Get(headerNameXSign)
	if !strings.HasPrefix(signature, signaturePrefix) {
		return errNoXSignHeader
	}

	// log.Printf("Signature at X-Hub = %s \n", signature)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	// We read the request body and now it's empty. We have to rewrite it for further reads.
	r.Body.Close() //nolint:errcheck
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// log.Printf("Request Body with r.Body = %v at Authorize check", body)

	validSignature, err := isValidSignature(signature, body)
	if err != nil {
		return fmt.Errorf("is valid signature: %w", err)
	}
	if !validSignature {
		return errInvalidXSignHeader
	}

	return nil
}

func signBody(body []byte) []byte {
	h := hmac.New(sha1.New, []byte(appSecret))
	h.Reset()
	h.Write(body)
	log.Println("signBody", h.Sum(nil))
	return h.Sum(nil)
}

func isValidSignature(signature string, body []byte) (bool, error) {
	actualSign, err := hex.DecodeString(signature[len(signaturePrefix):])
	if err != nil {
		log.Printf("decode string: %v", err)
		return false, fmt.Errorf("decode string: %w", err)
	}
	log.Println("----------------------------")
	log.Println("actualBody", actualSign)
	return hmac.Equal(signBody(body), actualSign), nil
}
