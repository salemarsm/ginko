package main

import (
	"net/http"
	"testing"
)

func TestSetBearer(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://example.test", nil)
	if err != nil {
		t.Fatal(err)
	}
	setBearer(req, " secret ")
	if got := req.Header.Get("Authorization"); got != "Bearer secret" {
		t.Fatalf("unexpected authorization header %q", got)
	}
}
