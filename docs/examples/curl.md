# curl examples

```bash
curl -X POST http://127.0.0.1:8787/api/context \
  -H 'content-type: application/json' \
  -d '{"query":"How should I answer?","subject":"botmaster","max_tokens":400}'
```

```bash
curl -X POST http://127.0.0.1:8787/api/suggest \
  -H 'content-type: application/json' \
  -d '{"subject":"botmaster","user_prompt":"I prefer direct Go examples."}'
```
