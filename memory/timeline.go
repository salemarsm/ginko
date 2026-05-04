package memory

import (
	"context"
	"database/sql"
)

// MemoryTimeline is an audit-focused lifecycle view for one memory.
type MemoryTimeline struct {
	Memory *Memory `json:"memory,omitempty"`
	Events []Event `json:"events"`
}

func (s *Store) MemoryTimeline(ctx context.Context, id string, limit int) (MemoryTimeline, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	var out MemoryTimeline
	if m, err := s.GetMemory(ctx, id); err == nil {
		out.Memory = &m
	} else if err != ErrNotFound {
		return MemoryTimeline{}, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, memory_id, kind, payload, source_kind, source_ref, created_at
		FROM events
		WHERE memory_id = ? OR payload LIKE ?
		ORDER BY created_at ASC
		LIMIT ?`, id, "%"+id+"%", limit)
	if err != nil {
		return MemoryTimeline{}, err
	}
	defer rows.Close()
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return MemoryTimeline{}, err
		}
		out.Events = append(out.Events, e)
	}
	if out.Memory == nil && len(out.Events) == 0 {
		return MemoryTimeline{}, ErrNotFound
	}
	return out, rows.Err()
}

func scanEvent(rows scanner) (Event, error) {
	var e Event
	var created string
	var memoryID sql.NullString
	if err := rows.Scan(&e.ID, &memoryID, &e.Kind, &e.Payload, &e.Source.Kind, &e.Source.Ref, &created); err != nil {
		return Event{}, err
	}
	if memoryID.Valid {
		e.MemoryID = &memoryID.String
	}
	t, err := parseTime(created)
	if err != nil {
		return Event{}, err
	}
	e.CreatedAt = t
	return e, nil
}
