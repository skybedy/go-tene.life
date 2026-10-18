package newsai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGenerateUsesStructuredResponsesAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("unexpected request: %s %q", r.Method, r.Header.Get("Authorization"))
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["model"] != "test-model" || body["store"] != false {
			t.Fatalf("unexpected model/store: %#v", body)
		}
		input, ok := body["input"].(string)
		if !ok || !strings.Contains(input, "ověřené podklady") {
			t.Fatalf("source material missing: %#v", body["input"])
		}
		text, ok := body["text"].(map[string]any)
		if !ok {
			t.Fatal("missing text configuration")
		}
		format := text["format"].(map[string]any)
		if format["type"] != "json_schema" || format["strict"] != true {
			t.Fatalf("structured output not strict: %#v", format)
		}
		response := map[string]any{
			"status": "completed",
			"output": []any{map[string]any{
				"type": "message",
				"content": []any{map[string]any{
					"type": "output_text",
					"text": `{"title":"Nová uzavírka","summary":"Krátký perex.","body":"## Omezení\n\nSilnice bude uzavřena.","category_slug":"doprava","seo_title":"Uzavírka na Tenerife","seo_description":"Aktuální informace o uzavírce.","editor_notes":"Ověřit čas."}`,
				}},
			}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	generator := &Generator{apiKey: "test-key", model: "test-model", endpoint: server.URL, client: &http.Client{Timeout: time.Second}}
	draft, err := generator.Generate(context.Background(), Input{
		SourceMaterial:  "ověřené podklady",
		EditorialBrief:  "stručně",
		CategoryChoices: []string{"doprava", "turistika"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if draft.Title != "Nová uzavírka" || draft.CategorySlug != "doprava" || !strings.Contains(draft.Body, "## Omezení") {
		t.Fatalf("unexpected draft: %+v", draft)
	}
}

func TestGenerateRejectsBadInputAndOutput(t *testing.T) {
	disabled := &Generator{}
	if _, err := disabled.Generate(context.Background(), Input{}); err != ErrNotConfigured {
		t.Fatalf("disabled generator: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"{\"title\":\"\"}"}]}]}`))
	}))
	defer server.Close()
	generator := &Generator{apiKey: "key", model: "model", endpoint: server.URL, client: server.Client()}
	if _, err := generator.Generate(context.Background(), Input{SourceMaterial: "text", CategoryChoices: []string{"doprava"}}); err != ErrInvalidOutput {
		t.Fatalf("invalid output accepted: %v", err)
	}
	if _, err := generator.Generate(context.Background(), Input{CategoryChoices: []string{"doprava"}}); err != ErrInvalidInput {
		t.Fatalf("empty input accepted: %v", err)
	}
}
