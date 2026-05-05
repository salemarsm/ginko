package memory

import (
	"fmt"
	"regexp"
)

// secretPattern describes a single sensitive-data pattern.
type secretPattern struct {
	name string
	re   *regexp.Regexp
}

var sensitivePatterns = []secretPattern{
	{"AWS access key", regexp.MustCompile(`AKIA[0-9A-Z]{16}`)},
	{"AWS secret key", regexp.MustCompile(`(?i)aws[_\-.]?secret[_\-.]?(?:access[_\-.]?)?key\s*[=:]\s*[A-Za-z0-9/+]{40}`)},
	{"private key header", regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----`)},
	{"GitHub token", regexp.MustCompile(`gh[ps]_[A-Za-z0-9]{36}`)},
	{"generic API key", regexp.MustCompile(`(?i)api[_\-.]?key\s*[=:]\s*['"]?[A-Za-z0-9\-_]{20,}['"]?`)},
	{"bearer token assignment", regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9\-_.~+/]{40,}`)},
	{"JWT token", regexp.MustCompile(`eyJ[A-Za-z0-9\-_]+\.eyJ[A-Za-z0-9\-_]+\.[A-Za-z0-9\-_.+/]*`)},
	{"password assignment", regexp.MustCompile(`(?i)password\s*[=:]\s*['"]?\S{8,}['"]?`)},
	{"database URL with credentials", regexp.MustCompile(`(?i)(?:postgres|mysql|mongodb|redis)://[^:@\s]+:[^@\s]+@`)},
}

// DetectSensitiveData returns an error if content contains patterns that look
// like secrets (API keys, private keys, passwords, tokens, etc.).
// Use <private>...</private> tags to wrap intentionally sensitive content
// that should be stripped before this check (handled upstream by StripPrivateTags).
func DetectSensitiveData(content string) error {
	for _, p := range sensitivePatterns {
		if p.re.MatchString(content) {
			return fmt.Errorf("memory content appears to contain sensitive data (%s); wrap in <private>...</private> to strip, or remove before saving", p.name)
		}
	}
	return nil
}

// SensitiveDataReport returns all detected pattern names without blocking.
// Useful for auditing existing memories.
func SensitiveDataReport(content string) []string {
	var found []string
	for _, p := range sensitivePatterns {
		if p.re.MatchString(content) {
			found = append(found, p.name)
		}
	}
	return found
}

