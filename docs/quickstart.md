# Quickstart

```bash
git clone https://github.com/salemarsm/llm-memory.git
cd llm-memory

go build -o bin/llm-memory ./cmd/llm-memory
go build -o bin/memserver ./cmd/memserver
go build -o bin/memmcp ./cmd/memmcp
go build -o bin/memctl ./cmd/memctl

bin/llm-memory init
bin/llm-memory doctor
bin/llm-memory ui
```

Open:

```txt
http://127.0.0.1:8787
```

Store a memory:

```bash
echo "The user prefers direct technical answers." \
  | bin/memctl -subject botmaster -scope global -type preference remember
```

Retrieve compact context:

```bash
bin/memctl -subject botmaster -scope global -max-tokens 400 context "How should I answer?"
```

Suggest learnings:

```bash
bin/memctl -subject botmaster suggest "I prefer Go examples and concise answers."
```
