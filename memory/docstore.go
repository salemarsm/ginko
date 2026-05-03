package memory

import (
	"context"
	"encoding/json"
	"time"
)

func (s *Store) UpsertDocument(ctx context.Context, d Document) (Document, error) {
	if d.ID == "" {
		d.ID = newID("doc")
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO documents(id, path, title, source_kind, source_ref, sha256, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET path=excluded.path, title=excluded.title, source_kind=excluded.source_kind, source_ref=excluded.source_ref, sha256=excluded.sha256`,
		d.ID, d.Path, d.Title, d.SourceKind, d.SourceRef, d.SHA256, formatTime(d.CreatedAt))
	return d, err
}

func (s *Store) UpsertChunk(ctx context.Context, c Chunk) (Chunk, error) {
	if c.ID == "" {
		c.ID = newID("chk")
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	if c.TokenCount <= 0 {
		c.TokenCount = EstimateTokens(c.Content)
	}
	if c.EmbeddingRefs == nil {
		c.EmbeddingRefs = EmbeddingRefs{}
	}
	embeds, _ := json.Marshal(c.EmbeddingRefs)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Chunk{}, err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO chunks(id, document_id, ordinal, heading_path, content, token_count, page_from, page_to, embedding_refs_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET document_id=excluded.document_id, ordinal=excluded.ordinal, heading_path=excluded.heading_path, content=excluded.content, token_count=excluded.token_count, page_from=excluded.page_from, page_to=excluded.page_to, embedding_refs_json=excluded.embedding_refs_json`,
		c.ID, c.DocumentID, c.Ordinal, c.HeadingPath, c.Content, c.TokenCount, nullableInt(c.PageFrom), nullableInt(c.PageTo), string(embeds), formatTime(c.CreatedAt))
	if err != nil {
		return Chunk{}, err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM chunks_fts WHERE id = ?`, c.ID); err != nil {
		return Chunk{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO chunks_fts(id, document_id, heading_path, content) VALUES (?, ?, ?, ?)`, c.ID, c.DocumentID, c.HeadingPath, c.Content); err != nil {
		return Chunk{}, err
	}
	return c, tx.Commit()
}

func nullableInt(i *int) any {
	if i == nil {
		return nil
	}
	return *i
}
