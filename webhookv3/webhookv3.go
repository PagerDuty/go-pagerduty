// Package webhookv3 provides functionality for working with V3 PagerDuty
// Webhooks, including signature verification and decoding.
package webhookv3

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ErrNoValidSignatures is returned when a webhook is not properly signed
// with the expected signature. When receiving this error, it is reccommended
// that the server return HTTP 403 to prevent redelivery.
var ErrNoValidSignatures = errors.New("invalid webhook signature")

// ErrMalformedHeader is returned when the *http.Request is missing the
// X-PagerDuty-Signature header. When receiving this error, it is recommended
// that the server return HTTP 400 to prevent redelivery.
var ErrMalformedHeader = errors.New("X-PagerDuty-Signature header is either missing or malformed")

// ErrMalformedBody is returned when the *http.Request body is either
// missing or malformed. When receiving this error, it's recommended that the
// server return HTTP 400 to prevent redelivery.
var ErrMalformedBody = errors.New("HTTP request body is either empty or malformed")

const (
	webhookSignaturePrefix = "v1="
	webhookSignatureHeader = "X-PagerDuty-Signature"
	webhookBodyReaderLimit = 2 * 1024 * 1024 // 2MB
)

// VerifySignature compares the provided signature of a PagerDuty v3 Webhook
// against the expected value and returns an ErrNoValidSignature error if the
// values do not match. This function may also return ErrMalformedHeader or
// ErrMalformedBody if the request appears to be malformed.
//
// See https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTkz-verifying-signatures for more details.
//
// This function will fail to read any HTTP request body that's 2MB or larger.
func VerifySignature(r *http.Request, secret string) error {
	h := r.Header.Get(webhookSignatureHeader)
	if len(h) == 0 {
		return ErrMalformedHeader
	}

	orb := r.Body

	b, err := io.ReadAll(io.LimitReader(r.Body, webhookBodyReaderLimit))
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	defer func() { _ = orb.Close() }()
	r.Body = io.NopCloser(bytes.NewReader(b))

	if len(b) == 0 {
		return ErrMalformedBody
	}

	sigs := extractPayloadSignatures(h)
	if len(sigs) == 0 {
		return ErrMalformedHeader
	}

	s := calculateSignature(b, secret)

	for _, sig := range sigs {
		if hmac.Equal(s, sig) {
			return nil
		}
	}

	return ErrNoValidSignatures
}

func extractPayloadSignatures(s string) [][]byte {
	var sigs [][]byte

	for _, sv := range strings.Split(s, ",") {
		// Ignore any signatures that are not the initial v1 version.
		if !strings.HasPrefix(sv, webhookSignaturePrefix) {
			continue
		}

		sig, err := hex.DecodeString(strings.TrimPrefix(sv, webhookSignaturePrefix))
		if err != nil {
			continue
		}

		sigs = append(sigs, sig)
	}

	return sigs
}

func calculateSignature(payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return mac.Sum(nil)
}
