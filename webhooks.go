package pagerduty

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	currentSignaturePrefix = "v1="
)

var (
	ErrNoValidSignature = errors.New("invalid webhook signature")
)

// VerifySignature compares the provided signature of a PagerDuty v3 Webhook
// against the expected value and returns an error if the values do not match.
//
// See https://developer.pagerduty.com/docs/webhooks/webhook-signatures/ for more details.
//
func VerifySignature(payload []byte, header, secret string) error {
	expectedSignature := calculateSignature(payload, secret)
	signatures := extractPayloadSignatures(header)

	for _, signature := range signatures {
		if hmac.Equal(expectedSignature, signature) {
			return nil
		}
	}

	return ErrNoValidSignature
}

func extractPayloadSignatures(signature string) (currentSignatures [][]byte) {
	signatureVersions := strings.Split(signature, ",")

	for _, signatureVersion := range signatureVersions {
		// Ignore any signatures that are not the initial v1 version.
		if !strings.HasPrefix(signatureVersion, currentSignaturePrefix) {
			continue
		}

		signature := strings.TrimPrefix(signatureVersion, currentSignaturePrefix)
		currentSignature, err := hex.DecodeString(signature)
		if err != nil {
			continue
		}

		currentSignatures = append(currentSignatures, currentSignature)
	}

	return currentSignatures
}

func calculateSignature(payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return mac.Sum(nil)
}
