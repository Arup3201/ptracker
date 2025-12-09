package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// https://github.com/nirasan/go-oauth-pkce-code-verifier/blob/master/verifier.go

const (
	DefaultLength = 32
	MinLength     = 32
	MaxLength     = 96
)

type CodeVerifier struct {
	Value string
}

func CreateCodeVerifier() (*CodeVerifier, error) {
	return CreateCodeVerifierWithLength(DefaultLength)
}

func CreateCodeVerifierWithLength(length int) (*CodeVerifier, error) {
	if length < MinLength || length > MaxLength {
		return nil, fmt.Errorf("invalid length: %v", length)
	}
	buf, err := randomBytes(length)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return CreateCodeVerifierFromBytes(buf)
}

func CreateCodeVerifierFromBytes(b []byte) (*CodeVerifier, error) {
	return &CodeVerifier{
		Value: encode(b),
	}, nil
}

func (v *CodeVerifier) CodeChallengePlain() string {
	return v.Value
}

func (v *CodeVerifier) CodeChallengeSHA256() string {
	h := sha256.New()
	h.Write([]byte(v.Value))
	return encode(h.Sum(nil))
}

func (v *CodeVerifier) String() string {
	return v.Value
}

func encode(msg []byte) string {
	encoded := base64.StdEncoding.EncodeToString(msg)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.ReplaceAll(encoded, "=", "")
	return encoded
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const csLen = byte(len(charset))

// https://tools.ietf.org/html/rfc7636#section-4.1)
func randomBytes(length int) ([]byte, error) {
	output := make([]byte, 0, length)
	for {
		buf := make([]byte, length)
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			return nil, fmt.Errorf("failed to read random bytes: %v", err)
		}
		for _, b := range buf {
			// Avoid bias by using a value range that's a multiple of 62
			if b < (csLen * 4) {
				output = append(output, charset[b%csLen])

				if len(output) == length {
					return output, nil
				}
			}
		}
	}
}
