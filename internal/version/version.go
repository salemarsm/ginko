package version

var (
	// Version is set by release builds with -ldflags "-X .../internal/version.Version=<tag>".
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

func String() string {
	return Version + " (commit " + Commit + ", built " + Date + ")"
}
