# Security model

## Design assumptions

- Local-first, single-user usage.
- Server binds to `127.0.0.1` by default.
- No secrets in memory content.
- No production multi-user isolation yet.

## HTTP API auth

Loopback-only binds (`127.0.0.1:8787`, `localhost:8787`, `[::1]:8787`) may run without auth for local development. Non-loopback binds (`0.0.0.0:8787`) fail config validation unless a bearer token is configured.

Prefer an environment variable over a literal token in config:

```json
{
  "server": {
    "addr": "0.0.0.0:8787",
    "auth_token_env": "GINKO_API_TOKEN"
  }
}
```

Clients include the token on all `/api` and `/api/v1` requests:

```http
Authorization: Bearer <token>
```

`GET /healthz` remains public for local health checks.

This is not production multi-user isolation; use it as single-user local/VPN/container protection.

## Sensitive-data detector

`memory_remember` and `POST /api/memories` run a pattern scan over memory content before persisting. Memories matching API key, token, or credential patterns are rejected with a validation error.

Patterns checked: `sk-...`, `ghp_...`, `AKIA...`, `-----BEGIN ... KEY-----`, and similar high-signal credential formats.

## Memory write policy

Agents should not store everything.

**Store:**
- explicit user preferences
- stable project facts
- architectural decisions
- corrections
- durable constraints
- long-lived tasks
- approved learnings

**Do not store:**
- transient chat context
- secrets or credentials
- sensitive personal data without explicit approval
- raw document chunks as memories
- uncertain inference as fact
- private data in shared/group contexts

## Private blocks

Content wrapped in `<private>...</private>` tags is stripped before persistence at all ingress points (MCP, HTTP API, store). A memory whose entire content is private is rejected.

## Network exposure

- Default bind: `127.0.0.1:8787` — accessible only from localhost.
- Docker: `-p 127.0.0.1:8787:8787` to restrict to localhost even in containers.
- Never expose directly to untrusted networks without TLS termination (e.g., nginx, Caddy, Tailscale).

---

## Backup and restore

The database is a standard SQLite 3 file at `~/.ginko/ginko.db` (or the path in `server.db`).

### Backup

**While the server is stopped:**

```bash
cp ~/.ginko/ginko.db ~/.ginko/ginko.db.bak
```

**While the server is running** (safe with WAL mode):

```bash
sqlite3 ~/.ginko/ginko.db ".backup $HOME/.ginko/ginko.db.bak"
```

Or with the SQLite online backup API via any SQLite client that supports it.

**Automated daily backup** (cron):

```cron
0 3 * * * sqlite3 ~/.ginko/ginko.db ".backup $HOME/.ginko/backups/ginko-$(date +\%Y\%m\%d).db"
```

### Restore

1. Stop the server: `pkill ginko` or `pkill memserver`.
2. Replace the database file:

```bash
cp ~/.ginko/ginko.db.bak ~/.ginko/ginko.db
```

3. Start the server: `ginko serve`.

### Verify integrity

```bash
sqlite3 ~/.ginko/ginko.db "PRAGMA integrity_check;"
# → ok
```

### Export memories as JSON

```bash
memctl -scope global list | jq . > memories-export.json
```

Or via the HTTP API:

```bash
curl http://127.0.0.1:8787/api/memories?limit=1000 > memories.json
```

### What is not backed up by the database file

- `~/.ginko/config.json` — back up separately if you have customized settings.
- Agent signal state — signals are ephemeral by design; leases expire automatically.
