package notify

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNtfySend(t *testing.T) {
	var gotMethod, gotPath, gotTitle, gotBody string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotTitle = r.Header.Get("Title")
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
	}))
	defer server.Close()

	n := NewNtfy(server.URL, "test-topic")

	err := n.Send(context.Background(), "Gold alert", "Crossed 3500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/test-topic" {
		t.Errorf("path = %q, want /test-topic", gotPath)
	}
	if gotTitle != "Gold alert" {
		t.Errorf("title = %q, want %q", gotTitle, "Gold alert")
	}
	if gotBody != "Crossed 3500" {
		t.Errorf("body = %q, want %q", gotBody, "Crossed 3500")
	}
}

func TestNtfySendError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := NewNtfy(server.URL, "test-topic")

	err := n.Send(context.Background(), "Gold alert", "Crossed 3500")
	if err == nil {
		t.Error("expected an error, got nil")
	}
}
