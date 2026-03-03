package pws

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// sanitizeLogText redacts the WEATHER_COM_API_KEY from any string.
// It handles plain text and URL-encoded values.
func sanitizeLogText(text string) string {
	apiKey := strings.TrimSpace(os.Getenv("WEATHER_COM_API_KEY"))
	if apiKey == "" || len(apiKey) < 8 {
		return text
	}

	redacted := "[REDACTED_API_KEY]"

	// 1. Plain API key
	text = strings.ReplaceAll(text, apiKey, redacted)

	// 2. URL-encoded API key
	encodedKey := url.QueryEscape(apiKey)
	if encodedKey != apiKey {
		text = strings.ReplaceAll(text, encodedKey, redacted)
	}

	return text
}

type redactedError struct {
	inner error
	msg   string
}

func (e *redactedError) Error() string {
	return e.msg
}

func (e *redactedError) Unwrap() error {
	return e.inner
}

// RedactError returns an error with sanitized message while preserving the original error chain.
func RedactError(err error) error {
	if err == nil {
		return nil
	}
	return &redactedError{
		inner: err,
		msg:   fmt.Sprintf("%s", sanitizeLogText(err.Error())),
	}
}
