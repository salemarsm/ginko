package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

func (s *Store) UpsertDocument(ctx context.Context, d Document) (Document, error) {
	if d.ID == "" {
		d.ID = newID("doc")
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO documents(id, ingestion_run_id, path, title, source_kind, source_ref, sha256, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET ingestion_run_id=excluded.ingestion_run_id, path=excluded.path, title=excluded.title, source_kind=excluded.source_kind, source_ref=excluded.source_ref, sha256=excluded.sha256`,
		d.ID, d.IngestionRunID, d.Path, d.Title, d.SourceKind, d.SourceRef, d.SHA256, formatTime(d.CreatedAt))
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

func (s *Store) UpsertIngestionRun(ctx context.Context, r IngestionRun) error {
	completed := ""
	if r.CompletedAt != nil {
		completed = formatTime(*r.CompletedAt)
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO ingestion_runs(id, source_path, recursive, parser, status, files_seen, documents_created, chunks_created, error, created_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET source_path=excluded.source_path, recursive=excluded.recursive, parser=excluded.parser, status=excluded.status, files_seen=excluded.files_seen, documents_created=excluded.documents_created, chunks_created=excluded.chunks_created, error=excluded.error, completed_at=excluded.completed_at`,
		r.ID, r.SourcePath, boolInt(r.Recursive), r.Parser, r.Status, r.FilesSeen, r.DocumentsCreated, r.ChunksCreated, r.Error, formatTime(r.CreatedAt), nullableEmptyTime(completed))
	return err
}

func (s *Store) ReplaceDocumentChunks(ctx context.Context, documentID string, contents []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM chunks_fts WHERE document_id = ?`, documentID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM chunks WHERE document_id = ?`, documentID); err != nil {
		return err
	}
	for i, content := range contents {
		c := Chunk{ID: newID("chk"), DocumentID: documentID, Ordinal: i, HeadingPath: "", Content: content, TokenCount: EstimateTokens(content), EmbeddingRefs: EmbeddingRefs{}, CreatedAt: time.Now().UTC()}
		embeds, _ := json.Marshal(c.EmbeddingRefs)
		if _, err := tx.ExecContext(ctx, `INSERT INTO chunks(id, document_id, ordinal, heading_path, content, token_count, page_from, page_to, embedding_refs_json, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, c.ID, c.DocumentID, c.Ordinal, c.HeadingPath, c.Content, c.TokenCount, nullableInt(c.PageFrom), nullableInt(c.PageTo), string(embeds), formatTime(c.CreatedAt)); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO chunks_fts(id, document_id, heading_path, content) VALUES (?, ?, ?, ?)`, c.ID, c.DocumentID, c.HeadingPath, c.Content); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListDocuments(ctx context.Context, limit int) ([]Document, error) {
	if limit <= 0 || limit > 1000 {
		limit = 200
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, COALESCE(ingestion_run_id, ''), path, title, source_kind, source_ref, sha256, created_at FROM documents ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Document
	for rows.Next() {
		var d Document
		var created string
		if err := rows.Scan(&d.ID, &d.IngestionRunID, &d.Path, &d.Title, &d.SourceKind, &d.SourceRef, &d.SHA256, &created); err != nil {
			return nil, err
		}
		d.CreatedAt, _ = parseTime(created)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *Store) ListChunks(ctx context.Context, documentID string) ([]Chunk, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, document_id, ordinal, heading_path, content, token_count, page_from, page_to, embedding_refs_json, created_at FROM chunks WHERE document_id = ? ORDER BY ordinal`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Chunk
	for rows.Next() {
		var c Chunk
		var pageFrom, pageTo sql.NullInt64
		var embeds, created string
		if err := rows.Scan(&c.ID, &c.DocumentID, &c.Ordinal, &c.HeadingPath, &c.Content, &c.TokenCount, &pageFrom, &pageTo, &embeds, &created); err != nil {
			return nil, err
		}
		if pageFrom.Valid {
			v := int(pageFrom.Int64)
			c.PageFrom = &v
		}
		if pageTo.Valid {
			v := int(pageTo.Int64)
			c.PageTo = &v
		}
		_ = json.Unmarshal([]byte(embeds), &c.EmbeddingRefs)
		c.CreatedAt, _ = parseTime(created)
		out = append(out, c)
	}
	return out, rows.Err()
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
func nullableEmptyTime(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func (s *Store) ListIngestionRuns(ctx context.Context, limit int) ([]IngestionRun, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, source_path, recursive, parser, status, files_seen, documents_created, chunks_created, error, created_at, completed_at FROM ingestion_runs ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []IngestionRun
	for rows.Next() {
		var r IngestionRun
		var recursive int
		var created string
		var completed sql.NullString
		if err := rows.Scan(&r.ID, &r.SourcePath, &recursive, &r.Parser, &r.Status, &r.FilesSeen, &r.DocumentsCreated, &r.ChunksCreated, &r.Error, &created, &completed); err != nil {
			return nil, err
		}
		r.Recursive = recursive == 1
		r.CreatedAt, _ = parseTime(created)
		if completed.Valid {
			if t, err := parseTime(completed.String); err == nil {
				r.CompletedAt = &t
			}
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) SearchChunks(ctx context.Context, req ChunkSearchRequest) ([]ChunkSearchResult, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	where := []string{"1=1"}
	args := []any{}
	join := ""
	orderBy := "c.created_at DESC, c.ordinal ASC"
	if fts := ftsQuery(req.Text); fts != "" {
		join = "JOIN chunks_fts ON chunks_fts.id = c.id"
		where = append(where, "chunks_fts MATCH ?")
		args = append(args, fts)
		orderBy = "bm25(chunks_fts) ASC, c.ordinal ASC"
	}
	if req.DocumentID != "" {
		where = append(where, "c.document_id = ?")
		args = append(args, req.DocumentID)
	}
	args = append(args, limit)
	query := `SELECT c.id, c.document_id, c.ordinal, c.heading_path, c.content, c.token_count, c.page_from, c.page_to, c.embedding_refs_json, c.created_at,
		d.id, COALESCE(d.ingestion_run_id, ''), d.path, d.title, d.source_kind, d.source_ref, d.sha256, d.created_at`
	if join != "" {
		query += `, bm25(chunks_fts)`
	} else {
		query += `, 0.0`
	}
	query += ` FROM chunks c ` + join + ` JOIN documents d ON d.id = c.document_id WHERE ` + strings.Join(where, " AND ") + ` ORDER BY ` + orderBy + ` LIMIT ?`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ChunkSearchResult
	for rows.Next() {
		var r ChunkSearchResult
		var pageFrom, pageTo sql.NullInt64
		var embeds, chunkCreated, docCreated string
		if err := rows.Scan(&r.Chunk.ID, &r.Chunk.DocumentID, &r.Chunk.Ordinal, &r.Chunk.HeadingPath, &r.Chunk.Content, &r.Chunk.TokenCount, &pageFrom, &pageTo, &embeds, &chunkCreated,
			&r.Document.ID, &r.Document.IngestionRunID, &r.Document.Path, &r.Document.Title, &r.Document.SourceKind, &r.Document.SourceRef, &r.Document.SHA256, &docCreated, &r.Score); err != nil {
			return nil, err
		}
		if pageFrom.Valid {
			v := int(pageFrom.Int64)
			r.Chunk.PageFrom = &v
		}
		if pageTo.Valid {
			v := int(pageTo.Int64)
			r.Chunk.PageTo = &v
		}
		_ = json.Unmarshal([]byte(embeds), &r.Chunk.EmbeddingRefs)
		r.Chunk.CreatedAt, _ = parseTime(chunkCreated)
		r.Document.CreatedAt, _ = parseTime(docCreated)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *Store) GetDocument(ctx context.Context, id string) (Document, error) {
	var d Document
	var created string
	err := s.db.QueryRowContext(ctx, `SELECT id, COALESCE(ingestion_run_id, ''), path, title, source_kind, source_ref, sha256, created_at FROM documents WHERE id = ?`, id).Scan(&d.ID, &d.IngestionRunID, &d.Path, &d.Title, &d.SourceKind, &d.SourceRef, &d.SHA256, &created)
	if err != nil {
		if err == sql.ErrNoRows {
			return Document{}, ErrNotFound
		}
		return Document{}, err
	}
	d.CreatedAt, _ = parseTime(created)
	return d, nil
}

func (s *Store) SuggestMemoriesFromDocument(ctx context.Context, req ChunkSuggestRequest) (ChunkSuggestResponse, error) {
	doc, err := s.GetDocument(ctx, req.DocumentID)
	if err != nil {
		return ChunkSuggestResponse{}, err
	}
	chunks, err := s.ListChunks(ctx, doc.ID)
	if err != nil {
		return ChunkSuggestResponse{}, err
	}
	limit := req.Limit
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	min := req.MinConfidence
	if min <= 0 {
		min = 0.65
	}
	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		subject = doc.Title
	}
	scope := req.Scope
	if scope == "" {
		scope = ScopeProject
	}
	seen := map[string]bool{}
	var candidates []MemoryCandidate
	for _, chunk := range chunks {
		for _, c := range suggestFromText(subject, scope, "chunk", chunk.Content) {
			if c.Memory.Confidence < min {
				continue
			}
			c.Memory.Source = Source{Kind: "chunk", Ref: doc.ID + ":" + chunk.ID}
			c.Memory.Tags = appendUnique(c.Memory.Tags, "rag", "evidence", "doc:"+doc.ID, "chunk:"+chunk.ID)
			key := strings.ToLower(string(c.Memory.Type) + "|" + c.Memory.Subject + "|" + c.Memory.Content)
			if seen[key] {
				continue
			}
			seen[key] = true
			candidates = append(candidates, c)
			if len(candidates) >= limit {
				break
			}
		}
		if len(candidates) >= limit {
			break
		}
	}
	resp := ChunkSuggestResponse{Document: doc, Candidates: candidates}
	if req.Store {
		for _, c := range candidates {
			m, err := s.UpsertMemory(ctx, c.Memory)
			if err != nil {
				return ChunkSuggestResponse{}, err
			}
			resp.Stored = append(resp.Stored, m)
		}
	}
	return resp, nil
}

func appendUnique(base []string, vals ...string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(base)+len(vals))
	for _, v := range append(base, vals...) {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}
