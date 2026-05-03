package memory

import "time"

// MemoryType classifies canonical memories. Keep this LLM-agnostic:
// models consume these records through APIs, but do not own the format.
type MemoryType string

const (
	TypePreference   MemoryType = "preference"
	TypeFact         MemoryType = "fact"
	TypeDecision     MemoryType = "decision"
	TypeTask         MemoryType = "task"
	TypeNote         MemoryType = "note"
	TypeRelationship MemoryType = "relationship"
)

// Scope controls where a memory may be retrieved.
type Scope string

const (
	ScopeGlobal  Scope = "global"
	ScopeProject Scope = "project"
	ScopeSession Scope = "session"
	ScopePrivate Scope = "private"
)

// Source preserves provenance. Memory without provenance is weak memory.
type Source struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

// EmbeddingRefs stores references to external/vector indexes.
// Embeddings are deliberately not canonical memory; they are disposable indexes.
type EmbeddingRefs map[string]string

// Memory is the canonical unit.
type Memory struct {
	ID            string        `json:"id"`
	Type          MemoryType    `json:"type"`
	Subject       string        `json:"subject"`
	Content       string        `json:"content"`
	Source        Source        `json:"source"`
	Scope         Scope         `json:"scope"`
	Confidence    float64       `json:"confidence"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	ValidFrom     *time.Time    `json:"valid_from,omitempty"`
	ValidUntil    *time.Time    `json:"valid_until,omitempty"`
	SupersedesID  *string       `json:"supersedes_id,omitempty"`
	SupersededBy  *string       `json:"superseded_by,omitempty"`
	Tags          []string      `json:"tags"`
	EmbeddingRefs EmbeddingRefs `json:"embedding_refs"`
}

// Event is append-only raw history. Canonical memories may be derived from it.
type Event struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Payload   string    `json:"payload"`
	Source    Source    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// Query is intentionally simple. More advanced systems can plug in BM25,
// vector search, graph traversal, or policy filters behind this API.
type Query struct {
	Text    string       `json:"text"`
	Types   []MemoryType `json:"types"`
	Scopes  []Scope      `json:"scopes"`
	Subject string       `json:"subject"`
	Tags    []string     `json:"tags"`
	Limit   int          `json:"limit"`
}
