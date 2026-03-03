package pws

import (
	"errors"
	"os"
	"testing"
)

func TestSanitizeLogText(t *testing.T) {
	testKey := "abcd-1234-xyz-789"
	t.Setenv("WEATHER_COM_API_KEY", testKey)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain_key",
			input:    "Your key is abcd-1234-xyz-789.",
			expected: "Your key is [REDACTED_API_KEY].",
		},
		{
			name:     "url_query_key",
			input:    "https://api.weather.com/v2/pws/observations/current?apiKey=abcd-1234-xyz-789&stationId=ICANARIA12",
			expected: "https://api.weather.com/v2/pws/observations/current?apiKey=[REDACTED_API_KEY]&stationId=ICANARIA12",
		},
		{
			name:     "json_payload",
			input:    `{"error":{"message":"Invalid apiKey 'abcd-1234-xyz-789' provided."}}`,
			expected: `{"error":{"message":"Invalid apiKey '[REDACTED_API_KEY]' provided."}}`,
		},
		{
			name:     "multiple_occurrences",
			input:    "Key1: abcd-1234-xyz-789, Key2: abcd-1234-xyz-789",
			expected: "Key1: [REDACTED_API_KEY], Key2: [REDACTED_API_KEY]",
		},
		{
			name:     "no_key_found",
			input:    "This message is safe.",
			expected: "This message is safe.",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeLogText(tt.input)
			if result != tt.expected {
				t.Fatalf("unexpected sanitize result\nexpected: %q\ngot:      %q", tt.expected, result)
			}
		})
	}
}

func TestSanitizeLogText_UrlEncoded(t *testing.T) {
	testKey := "key with spaces"
	t.Setenv("WEATHER_COM_API_KEY", testKey)

	input := "api.php?k=key+with+spaces&id=1"
	expected := "api.php?k=[REDACTED_API_KEY]&id=1"

	result := sanitizeLogText(input)
	if result != expected {
		t.Fatalf("unexpected sanitize result\nexpected: %q\ngot:      %q", expected, result)
	}
}

func TestSanitizeLogText_EmptyEnv(t *testing.T) {
	os.Unsetenv("WEATHER_COM_API_KEY")

	input := "Some abcd-1234 key"
	expected := input

	result := sanitizeLogText(input)
	if result != expected {
		t.Fatalf("unexpected sanitize result\nexpected: %q\ngot:      %q", expected, result)
	}
}

func TestRedactError_PreservesErrorChain(t *testing.T) {
	t.Setenv("WEATHER_COM_API_KEY", "abcd-1234-xyz-789")
	base := errors.New("request failed: apiKey=abcd-1234-xyz-789")

	err := RedactError(base)
	if err == nil {
		t.Fatalf("expected redacted error")
	}
	if !errors.Is(err, base) {
		t.Fatalf("expected wrapped error chain to contain base error")
	}
	if got := err.Error(); got == base.Error() {
		t.Fatalf("expected sanitized error text, got unchanged: %q", got)
	}
	if got := err.Error(); got != "request failed: apiKey=[REDACTED_API_KEY]" {
		t.Fatalf("unexpected redacted error text: %q", got)
	}
}
