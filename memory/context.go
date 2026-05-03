package memory

import (
	"context"
	"fmt"
	"strings"
)

const defaultContextTokenBudget = 1200

type ContextRequest struct {
	Query       string       `json:"query"`
	Subject     string       `json:"subject"`
	Types       []MemoryType `json:"types"`
	Scopes      []Scope      `json:"scopes"`
	Tags        []string     `json:"tags"`
	MaxTokens   int          `json:"max_tokens"`
	MaxMemories int          `json:"max_memories"`
}

type ContextResponse struct {
	Context         string   `json:"context"`
	Items           []Memory `json:"items"`
	EstimatedTokens int      `json:"estimated_tokens"`
	BudgetTokens    int      `json:"budget_tokens"`
	Truncated       bool     `json:"truncated"`
}

func (s *Store) BuildContext(ctx context.Context, req ContextRequest) (ContextResponse, error) {
	budget := req.MaxTokens
	if budget <= 0 {
		budget = defaultContextTokenBudget
	}
	if budget > 8000 {
		budget = 8000
	}
	limit := req.MaxMemories
	if limit <= 0 {
		limit = 12
	}
	if limit > 50 {
		limit = 50
	}

	items, err := s.Search(ctx, Query{
		Text:    req.Query,
		Types:   req.Types,
		Scopes:  req.Scopes,
		Subject: req.Subject,
		Tags:    req.Tags,
		Limit:   limit,
	})
	if err != nil {
		return ContextResponse{}, err
	}

	var b strings.Builder
	used := 0
	selected := make([]Memory, 0, len(items))
	truncated := false
	for _, m := range items {
		line := formatContextMemory(m)
		cost := EstimateTokens(line)
		if used+cost > budget {
			truncated = true
			break
		}
		b.WriteString(line)
		b.WriteByte('\n')
		used += cost
		selected = append(selected, m)
	}

	return ContextResponse{
		Context:         strings.TrimSpace(b.String()),
		Items:           selected,
		EstimatedTokens: used,
		BudgetTokens:    budget,
		Truncated:       truncated,
	}, nil
}

func formatContextMemory(m Memory) string {
	tags := ""
	if len(m.Tags) > 0 {
		tags = " tags=" + strings.Join(m.Tags, ",")
	}
	return fmt.Sprintf("- [%s/%s conf=%.2f src=%s:%s%s] %s", m.Type, m.Scope, m.Confidence, m.Source.Kind, m.Source.Ref, tags, compactWhitespace(m.Content))
}

func EstimateTokens(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	// Cheap, deterministic approximation good enough for budgeting before a provider-specific tokenizer exists.
	return (len([]rune(s)) / 4) + 1
}

func compactWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
