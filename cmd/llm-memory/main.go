package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/salemarsm/llm-memory/config"
	"github.com/salemarsm/llm-memory/internal/version"
)

const defaultDirName = ".ginko"

func main() {
	_ = config.MaybeMigrateLegacyDataDir()
	home, _ := os.UserHomeDir()
	defaultHome := filepath.Join(home, defaultDirName)
	homeDir := flag.String("home", envDefault("LLM_MEMORY_HOME", defaultHome), "llm-memory home directory")
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}

	cmd := flag.Arg(0)
	args := flag.Args()[1:]
	switch cmd {
	case "init":
		must(initProject(*homeDir))
	case "doctor":
		must(doctor(*homeDir))
	case "token":
		must(tokenCommand(*homeDir, args))
	case "setup":
		must(setupCommand(*homeDir, args))
	case "version":
		fmt.Println("llm-memory", version.String())
	case "paths":
		printPaths(*homeDir)
	case "mcp-config":
		must(printMCPConfig(*homeDir))
	case "install-mcp":
		must(installMCP(*homeDir, args))
	case "ui":
		must(runMemServer(*homeDir))
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

func initProject(home string) error {
	if err := os.MkdirAll(home, 0o755); err != nil {
		return err
	}
	cfgPath := configPath(home)
	if _, err := os.Stat(cfgPath); errors.Is(err, os.ErrNotExist) {
		cfg := config.Default()
		cfg.Database.Path = dbPath(home)
		b, _ := json.MarshalIndent(cfg, "", "  ")
		if err := os.WriteFile(cfgPath, append(b, '\n'), 0o644); err != nil {
			return err
		}
	}
	fmt.Println("✓ initialized", home)
	fmt.Println("config:", cfgPath)
	fmt.Println("database:", dbPath(home))
	return nil
}

func doctor(home string) error {
	fmt.Println("llm-memory doctor")
	fmt.Println("home:", home)
	cfgPath := configPath(home)
	if _, err := os.Stat(home); err != nil {
		fmt.Printf("✗ home dir: %s\n  fix: llm-memory -home %q init\n", err, home)
	} else {
		fmt.Println("✓ home dir:", home)
	}
	if _, err := os.Stat(cfgPath); err != nil {
		fmt.Printf("✗ config: %s\n  fix: llm-memory -home %q init\n", err, home)
	} else {
		fmt.Println("✓ config:", cfgPath)
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Printf("✗ config valid: %s\n  fix: edit %s or run llm-memory -home %q init on a clean home\n", err, cfgPath, home)
		cfg = config.Default()
	} else {
		fmt.Println("✓ config valid")
	}
	if config.IsLoopbackAddr(cfg.Server.Addr) {
		fmt.Println("✓ auth policy: loopback no-auth allowed")
	} else if _, ok := cfg.Server.BearerToken(); ok {
		fmt.Println("✓ auth policy: non-loopback bind has bearer token")
	} else {
		fmt.Println("✗ auth policy: non-loopback bind requires server.auth_token_env or server.auth_token")
	}
	for _, bin := range []string{"memserver", "memmcp", "memctl"} {
		p, err := findSibling(bin)
		if err != nil {
			fmt.Printf("✗ %s: %s\n  fix: run make build or install release artifacts\n", bin, err)
		} else {
			fmt.Printf("✓ %s: %s\n", bin, p)
		}
	}
	if canListen(cfg.Server.Addr) {
		fmt.Println("✓ port", cfg.Server.Addr, "available")
	} else {
		fmt.Println("! port", cfg.Server.Addr, "unavailable or already in use")
		fmt.Println("  fix: stop the running service or change server.addr in", cfgPath)
	}
	return nil
}

func tokenCommand(home string, args []string) error {
	if len(args) < 1 {
		return errors.New("token requires subcommand: create, list, revoke")
	}
	if err := initProject(home); err != nil {
		return err
	}
	cfgPath := configPath(home)
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	switch args[0] {
	case "create":
		token, err := randomToken()
		if err != nil {
			return err
		}
		cfg.Server.AuthToken = token
		cfg.Server.AuthTokenEnv = ""
		if err := writeConfig(cfgPath, cfg); err != nil {
			return err
		}
		fmt.Println(token)
		fmt.Fprintln(os.Stderr, "✓ wrote server.auth_token to", cfgPath)
	case "list":
		if cfg.Server.AuthToken != "" {
			fmt.Println("auth_token: configured")
		} else {
			fmt.Println("auth_token: not configured")
		}
		if cfg.Server.AuthTokenEnv != "" {
			_, ok := cfg.Server.BearerToken()
			fmt.Printf("auth_token_env: %s (set=%v)\n", cfg.Server.AuthTokenEnv, ok)
		} else {
			fmt.Println("auth_token_env: not configured")
		}
	case "revoke":
		cfg.Server.AuthToken = ""
		cfg.Server.AuthTokenEnv = ""
		if err := writeConfig(cfgPath, cfg); err != nil {
			return err
		}
		fmt.Fprintln(os.Stderr, "✓ cleared server auth token config in", cfgPath)
	default:
		return fmt.Errorf("unknown token subcommand %q", args[0])
	}
	return nil
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func writeConfig(path string, cfg config.Config) error {
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o600)
}

func printPaths(home string) {
	fmt.Println("home=", home)
	fmt.Println("config=", configPath(home))
	fmt.Println("db=", dbPath(home))
}

func printMCPConfig(home string) error {
	memmcp, err := findSibling("memmcp")
	if err != nil {
		return err
	}
	cfg := map[string]any{
		"command": memmcp,
		"args":    []string{"-db", dbPath(home)},
	}
	b, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(b))
	return nil
}

func setupCommand(home string, args []string) error {
	if len(args) < 1 {
		return errors.New("setup requires target: claude-code")
	}
	fs := flag.NewFlagSet("setup", flag.ContinueOnError)
	dryRun := fs.Bool("dry-run", false, "show changes without writing")
	local := fs.Bool("local", false, "write .claude/settings.json in current directory")
	configFile := fs.String("config", "", "explicit Claude Code settings.json path")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	switch args[0] {
	case "claude-code":
		return setupClaudeCode(home, *dryRun, *local, *configFile)
	default:
		return fmt.Errorf("unknown setup target %q", args[0])
	}
}

func setupClaudeCode(home string, dryRun, local bool, explicitPath string) error {
	if err := initProject(home); err != nil {
		return err
	}
	memmcp, err := findSibling("memmcp")
	if err != nil {
		return err
	}
	path, err := claudeSettingsPath(local, explicitPath)
	if err != nil {
		return err
	}
	original, merged, err := mergeClaudeSettings(path, map[string]any{
		"command": memmcp,
		"args":    []string{"-db", dbPath(home)},
	})
	if err != nil {
		return err
	}
	fmt.Println("target:", path)
	fmt.Println("mcp server: ginko")
	if dryRun {
		fmt.Println("dry-run: no files written")
		fmt.Println(string(merged))
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if original != nil {
		backup := path + ".bak"
		if err := os.WriteFile(backup, original, 0o600); err != nil {
			return err
		}
		fmt.Println("backup:", backup)
	}
	if err := os.WriteFile(path, merged, 0o600); err != nil {
		return err
	}
	fmt.Println("✓ configured Claude Code MCP server 'ginko'")
	fmt.Println("settings:", path)
	return nil
}

func claudeSettingsPath(local bool, explicitPath string) (string, error) {
	if strings.TrimSpace(explicitPath) != "" {
		return explicitPath, nil
	}
	if local {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(cwd, ".claude", "settings.json"), nil
	}
	cwd, err := os.Getwd()
	if err == nil {
		project := filepath.Join(cwd, ".claude", "settings.json")
		if _, statErr := os.Stat(project); statErr == nil {
			return project, nil
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

func mergeClaudeSettings(path string, server map[string]any) ([]byte, []byte, error) {
	settings := map[string]any{}
	var original []byte
	if b, err := os.ReadFile(path); err == nil {
		original = b
		if len(strings.TrimSpace(string(b))) > 0 {
			if err := json.Unmarshal(b, &settings); err != nil {
				return nil, nil, fmt.Errorf("parse %s: %w", path, err)
			}
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, nil, err
	}
	servers, _ := settings["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}
	servers["ginko"] = server
	settings["mcpServers"] = servers
	merged, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	return original, append(merged, '\n'), nil
}

func installMCP(home string, args []string) error {
	if len(args) < 1 {
		return errors.New("install-mcp requires target: claude-code, codex, openclaw, or print")
	}
	if err := initProject(home); err != nil {
		return err
	}
	target := args[0]
	memmcp, err := findSibling("memmcp")
	if err != nil {
		return err
	}
	snippet := map[string]any{"command": memmcp, "args": []string{"-db", dbPath(home)}}
	b, _ := json.MarshalIndent(snippet, "", "  ")
	switch target {
	case "print":
		fmt.Println(string(b))
	case "claude-code", "codex", "openclaw":
		out := filepath.Join(home, "mcp-"+target+".json")
		if err := os.WriteFile(out, append(b, '\n'), 0o644); err != nil {
			return err
		}
		fmt.Println("✓ wrote MCP config snippet:", out)
		fmt.Println("Add this MCP server to", target, "configuration.")
		fmt.Println("Bootstrap instruction:")
		fmt.Println(bootstrapInstruction())
	default:
		return fmt.Errorf("unknown MCP target %q", target)
	}
	return nil
}

func runMemServer(home string) error {
	if err := initProject(home); err != nil {
		return err
	}
	memserver, err := findSibling("memserver")
	if err != nil {
		return err
	}
	cmd := exec.Command(memserver, "-config", configPath(home))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func bootstrapInstruction() string {
	return strings.TrimSpace(`Before answering, silently call memory_context with the user request, subject, relevant scopes, and max_tokens <= 1200.
Do not mention memory unless asked.
After answering, call memory_suggest with the user prompt, assistant response, and a concise LLM inference about durable learnings.
Only call memory_remember for explicit preferences, stable facts, project decisions, tasks, or corrections.
Ask before storing sensitive, private, or uncertain information.
Prefer compact memories over raw document chunks.`)
}

func configPath(home string) string { return filepath.Join(home, "config.json") }
func dbPath(home string) string     { return filepath.Join(home, "ginko.db") }

func findSibling(name string) (string, error) {
	exe, err := os.Executable()
	if err == nil {
		candidate := filepath.Join(filepath.Dir(exe), exeName(name))
		if isExecutable(candidate) {
			return candidate, nil
		}
	}
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("%s not found next to llm-memory or in PATH", name)
}

func exeName(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func isExecutable(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir() && st.Mode()&0o111 != 0
}

func canListen(addr string) bool {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

func envDefault(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `llm-memory [flags] <command>

Commands:
  init                    create ~/.ginko/config.json and database path
  doctor                  check binaries, config, auth policy, and port
  token create|list|revoke manage local API bearer token config
  setup claude-code       configure Claude Code MCP server (use --dry-run first)
  version                 print version, commit, and build date
  paths                   print effective paths
  mcp-config              print MCP server JSON snippet
  install-mcp <target>    write MCP snippet for claude-code, codex, openclaw, or print
  ui                      run memserver with local config

Flags:
  -home DIR               default ~/.ginko or LLM_MEMORY_HOME`)
}
