package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/salemarsm/llm-memory/config"
	"github.com/salemarsm/llm-memory/memory"
	"github.com/salemarsm/llm-memory/server"
)

func main() {
	configPath := flag.String("config", "", "path to JSON config")
	writeConfig := flag.String("write-config", "", "write default JSON config and exit")
	flag.Parse()

	if *writeConfig != "" {
		if err := config.WriteDefault(*writeConfig); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, "wrote", *writeConfig)
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	store, err := memory.Open(cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	srv := server.New(store, cfg)
	log.Printf("llm-memory listening on http://%s", cfg.Server.Addr)
	log.Printf("database=%s llm=%s/%s embedding=%s/%s", cfg.Database.Path, cfg.LLM.Provider, cfg.LLM.Model, cfg.Embedding.Provider, cfg.Embedding.Model)
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, srv.Handler()))
}
