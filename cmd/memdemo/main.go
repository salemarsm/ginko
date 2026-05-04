package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/salemarsm/llm-memory/config"
	"github.com/salemarsm/llm-memory/memory"
)

func main() {
	_ = config.MaybeMigrateLegacyDataDir()
	ctx := context.Background()
	path := "memory.db"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	store, err := memory.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	pref, err := store.UpsertMemory(ctx, memory.Memory{
		Type:       memory.TypePreference,
		Subject:    "botmaster",
		Content:    "Prefere respostas diretas, técnicas e sem enrolação.",
		Source:     memory.Source{Kind: "conversation", Ref: "bootstrap/2026-05-03T13:29Z"},
		Scope:      memory.ScopeGlobal,
		Confidence: 0.95,
		Tags:       []string{"style", "preference"},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = store.AppendEvent(ctx, memory.Event{
		Kind:    "memory.created",
		Payload: pref.ID,
		Source:  memory.Source{Kind: "system", Ref: "memdemo"},
	})
	if err != nil {
		log.Fatal(err)
	}

	items, err := store.Search(ctx, memory.Query{
		Text:    "respostas diretas",
		Subject: "botmaster",
		Scopes:  []memory.Scope{memory.ScopeGlobal},
		Limit:   10,
	})
	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(items, "", "  ")
	fmt.Println(string(b))
}
