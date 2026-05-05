package main

import (
	"testing"

	"github.com/salemarsm/ginko/memory"
)

func openTestStore(t *testing.T) *memory.Store {
	t.Helper()
	store, err := memory.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}
