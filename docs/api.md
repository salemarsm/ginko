# HTTP API

See [openapi.yaml](openapi.yaml) for machine-readable API documentation.

Core endpoints are available under both `/api/...` and `/api/v1/...` during v0.x. Prefer `/api/v1/...` for new integrations.

Core endpoints:

- `POST /api/context`
- `POST /api/suggest`
- `POST /api/memories`
- `POST /api/search`
- `POST /api/supersede/{id}`
- `DELETE /api/memories/{id}`
- `GET /api/events`
- `GET /api/config`
- `GET /healthz`


Memory suggestion details: [Suggestion engine](suggestion-engine.md).
