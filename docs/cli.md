# CLI

## Ginko

```bash
bin/ginko init
bin/ginko doctor
bin/ginko paths
bin/ginko mcp-config
bin/ginko install-mcp claude-code
bin/ginko install-mcp codex
bin/ginko install-mcp openclaw
bin/ginko ui
```

## memctl

```bash
bin/memctl -subject botmaster search "direct answers"
bin/memctl -subject botmaster -max-tokens 800 context "How should I answer?"
bin/memctl -subject botmaster -type preference remember "User prefers direct answers."
bin/memctl -subject botmaster suggest "I prefer Go examples."
```

## Auth tokens

```sh
bin/ginko token create   # prints and stores a local bearer token
bin/ginko token list     # shows whether auth is configured; never prints secrets
bin/ginko token revoke   # clears token config

GINKO_API_TOKEN=<token> bin/memctl search "query"
bin/memctl -token <token> search "query"
```
