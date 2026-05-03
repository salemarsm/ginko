# Security model

Current default assumption:

- local-only usage
- bind to `127.0.0.1`
- no public network exposure
- no secrets in memory content
- no production multi-user isolation yet

Do **not** expose the HTTP server to untrusted networks until API token/local auth is implemented.

## Status

Current security posture is suitable for local experimentation only. API token support is planned for v0.2 and should be implemented before any remote/container/VPS exposure.

## Memory write policy

Agents should not store everything.

Store:

- explicit user preferences
- stable project facts
- architectural decisions
- corrections
- durable constraints
- long-lived tasks
- approved learnings

Do not store:

- transient chat context
- secrets or credentials
- sensitive personal data without explicit approval
- raw document chunks as memories
- uncertain inference as fact
- private data in shared/group contexts
