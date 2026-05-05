package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/salemarsm/llm-memory/config"
	"github.com/salemarsm/llm-memory/memory"
	"github.com/salemarsm/llm-memory/server"
)

func TestMemserverSmoke(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	store, err := memory.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	cfg := config.Default()
	cfg.Server.Addr = addr
	h := server.New(store, cfg).Handler()

	httpSrv := &http.Server{Addr: addr, Handler: h}
	go httpSrv.ListenAndServe() //nolint:errcheck

	base := "http://" + addr
	waitReady(t, base+"/healthz")

	t.Run("healthz", func(t *testing.T) {
		resp := mustGET(t, base+"/healthz")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
	})

	t.Run("config", func(t *testing.T) {
		resp := mustGET(t, base+"/api/config")
		var m map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if _, ok := m["server"]; !ok {
			t.Fatal("missing server key in /api/config")
		}
	})

	t.Run("memories_empty", func(t *testing.T) {
		resp := mustGET(t, base+"/api/memories")
		var items []any
		if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if len(items) != 0 {
			t.Fatalf("expected empty store, got %d items", len(items))
		}
	})

	t.Run("browse_home", func(t *testing.T) {
		resp := mustGET(t, base+"/api/browse")
		var m map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if _, ok := m["path"]; !ok {
			t.Fatal("missing path key in /api/browse")
		}
	})
}

func waitReady(t *testing.T, url string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:noctx
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("server at %s not ready after 3s", url)
}

func mustGET(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(fmt.Sprintf("%s", url)) //nolint:noctx
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	return resp
}
