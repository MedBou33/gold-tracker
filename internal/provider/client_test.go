package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"currency":"EUR","name":"Gold","price":3500.5,"symbol":"XAU","updatedAt":"2026-10-20T14:00:00Z"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	quote, err := client.GetQuote(context.Background(), "XAU", "EUR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if quote.Price != 3500.5 {
		t.Errorf("price = %v, want 3500.5", quote.Price)
	}
	if quote.Currency != "EUR" {
		t.Errorf("currency = %q, want EUR", quote.Currency)
	}
}

func TestGetQuoteErrors(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
	}{
		{"server error", http.StatusInternalServerError, "oops"},
		{"not found", http.StatusNotFound, "{}"},
		{"broken json", http.StatusOK, "not json"},
		{"zero price", http.StatusOK, `{"currency":"EUR","price":0}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := NewClient(server.URL)
			_, err := client.GetQuote(context.Background(), "XAU", "EUR")
			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}
