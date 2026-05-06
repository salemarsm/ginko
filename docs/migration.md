# Schema migration and versioning

## How migrations work

Ginko uses a simple append-only migration system backed by the `schema_migrations` table in SQLite.

On every `Store.Open`, the store runs all pending migrations in order. Each migration is a SQL statement identified by an integer version. If the migration has already been applied (its version exists in `schema_migrations`), it is skipped. Migrations are never rolled back automatically.

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL
);
```

The current schema version is reported by `ginko doctor` and `/api/config`.

## Checking the schema version

```bash
ginko doctor
# ✓ schema version: 1
```

Or via the API:

```bash
curl http://127.0.0.1:8787/api/config | jq .schema_version
```

## Forward migrations only

Ginko does not support rollback. The canonical upgrade path is:

1. Stop the server.
2. Back up the database (see [security.md](security.md#backup-and-restore)).
3. Install the new binary.
4. Start the server — migrations run automatically on startup.

## Adding a migration (contributors)

Append a new entry to the `migrations` slice in `memory/store.go`:

```go
{version: N, sql: `ALTER TABLE memories ADD COLUMN new_field TEXT;`},
```

Rules:
- `version` must be strictly greater than all previous entries.
- Each migration must be idempotent where possible (`IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`).
- Never modify an already-applied migration — create a new one instead.
- Test with both a fresh database and an existing database at the previous version.

## Compatibility guarantees

- **v0.x**: no stability guarantee. Schemas may change between minor releases. Always back up before upgrading.
- **v1.0+**: additive changes only (new tables, new nullable columns). Breaking schema changes require a new major version.

## SQLite file portability

The database is a standard SQLite 3 file. Any SQLite tool can read it:

```bash
sqlite3 ~/.ginko/ginko.db ".tables"
sqlite3 ~/.ginko/ginko.db "SELECT version, applied_at FROM schema_migrations ORDER BY version;"
```
