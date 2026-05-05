# Quickstart

## Install from release

Download the latest `ginko` archive for your OS/architecture from GitHub Releases, unpack it, and put the binaries on `PATH`.

The release archive contains:

- `ginko`
- `ginko-admin` (support binary used by `ginko`)
- `memctl`
- `memmcp`
- `memserver`

Verify:

```bash
ginko version
memmcp -version
memctl -version
memserver -version
```

## Build from source

```bash
git clone https://github.com/salemarsm/ginko.git
cd ginko

go build -o bin/ginko ./cmd/ginko
go build -o bin/ginko-admin ./cmd/ginko-admin
go build -o bin/memserver ./cmd/memserver
go build -o bin/memmcp ./cmd/memmcp
go build -o bin/memctl ./cmd/memctl

bin/ginko init
bin/ginko doctor
bin/ginko ui
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
