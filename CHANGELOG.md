# Changelog

All notable changes to this project are documented here.

## [Unreleased]

### Changed

- Project renamed definitively to **Ginko**.
- Go module path changed to `github.com/salemarsm/ginko`.
- Public docs, release metadata, install script, white paper paths, and examples now use Ginko naming.
- Primary CLI is `ginko`; lower-level support binaries remain `memctl`, `memmcp`, `memserver`, and internal helper `ginko-admin`.
- Default data directory is `~/.ginko/` and default database is `~/.ginko/ginko.db`.

### Added

- `ginko` umbrella CLI and Claude Code plugin (`plugin/claude-code/`). The plugin includes an MCP server entry, a Memory Protocol skill, SessionStart and PreCompact hooks, and slash commands `/ginko:save` and `/ginko:recall`.
- `ginko setup claude-code` and `ginko setup openclaw` for agent MCP setup.
- Marketplace manifest at `.claude-plugin/marketplace.json` exposing the `ginko` plugin. Install with:
      /plugin marketplace add salemarsm/ginko
      /plugin install ginko

### Compatibility

- Existing `~/.ginko/` data is preserved. Ginko no longer documents or emits legacy project names.
