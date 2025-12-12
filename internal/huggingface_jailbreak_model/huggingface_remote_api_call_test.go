package huggingface_jailbreak_model

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/safetorun/PromptDefender/wall"
)

func TestCallRemoteApi_Injection(t *testing.T) {
	// Mock server returning INJECTION with high score
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		response := `[[{"label": "INJECTION", "score": 0.99}]]`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
	defer server.Close()

	caller := NewRemoteApiCaller("test-token")
	// Override URL for testing
	caller.huggingfaceUrl = server.URL

	level, err := caller.CallRemoteApi("ignore me")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if level != wall.VeryClose && level != wall.ExactMatch {
		t.Errorf("Expected VeryClose or ExactMatch, got %v", level)
	}
}

func TestCallRemoteApi_Safe(t *testing.T) {
	// Mock server returning INJECTION with low score (safe)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `[[{"label": "INJECTION", "score": 0.01}]]`
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))
	defer server.Close()

	caller := NewRemoteApiCaller("test-token")
	caller.huggingfaceUrl = server.URL

	level, err := caller.CallRemoteApi("hello world")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if level != wall.NoMatch {
		t.Errorf("Expected NoMatch, got %v", level)
	}
}
