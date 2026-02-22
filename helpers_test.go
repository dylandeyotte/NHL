package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlayingToday(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"games":[{"gameDate":"2026-02-20"}, {"gameDate":"2026-02-21"}, {"gameDate":"2026-02-22"}, {"gameDate":"2026-02-23"}]}`))
	}))

	defer server.Close()

	client := server.Client()

	test, err := playingToday("VGK", server.URL, client)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if test != true {
		t.Fatalf("expected true, got %v", test)
	}

}
