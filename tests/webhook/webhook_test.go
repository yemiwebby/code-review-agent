package webhook_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yemiwebby/code-review-agent/internal/webhook"
)

func TestWebhookHandler(t *testing.T) {
	t.Skip("Temporarily skipping test for this handler to avoid GitHub Post calls")

	payload := []byte(`{
	    "action": "opened",
		"pull_request": {
		  "number": 1,
		  "title": "Update main.go",
		  "user": {
		    "login": "yemiwebby"
		  }
		},
		"repository": {
		    "full_name": "yemiwebby/test-ai-go-pr"
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	webhook.Handle(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rec.Result().StatusCode)
	}

}
