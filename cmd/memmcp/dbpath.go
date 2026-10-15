package main

import (
	"strings"

	"github.com/salemarsm/ginko/config"
)

// resolveDBPath decide onde o banco vive quando nem -db nem GINKO_DB foram dados.
//
// O memserver, o `ginko init` e a GUI resolvem o banco pelo config do ginko
// (config.Database.Path, tipicamente ~/.ginko/ginko.db). O memmcp ignorava isso
// e caia em "./memory.db", relativo ao diretorio de trabalho -- que num servidor
// MCP lancado pelo cliente e o diretorio DO CLIENTE, nao o do repositorio.
//
// O efeito era silencioso e grave: o agente gravava memoria num banco que o
// usuario nunca inspecionava, enquanto a GUI e o memctl liam outro.
func resolveDBPath(explicit string) string {
	if strings.TrimSpace(explicit) != "" {
		return explicit
	}
	if cfg, err := config.Load(config.DefaultConfigPath()); err == nil {
		if p := strings.TrimSpace(cfg.Database.Path); p != "" {
			return p
		}
	}
	return config.DefaultDBPath()
}
