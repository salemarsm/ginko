package memory

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type IngestRequest struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

type IngestResponse struct {
	Run       IngestionRun `json:"run"`
	Documents []Document   `json:"documents"`
	Chunks    []Chunk      `json:"chunks"`
	Skipped   []string     `json:"skipped"`
}

type IngestionRun struct {
	ID               string     `json:"id"`
	SourcePath       string     `json:"source_path"`
	Recursive        bool       `json:"recursive"`
	Parser           string     `json:"parser"`
	Status           string     `json:"status"`
	FilesSeen        int        `json:"files_seen"`
	DocumentsCreated int        `json:"documents_created"`
	ChunksCreated    int        `json:"chunks_created"`
	Error            string     `json:"error,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
}

var textIngestExts = map[string]bool{
	".txt": true, ".md": true, ".markdown": true, ".html": true, ".htm": true,
	".json": true, ".jsonl": true, ".csv": true, ".tsv": true, ".tex": true,
}

func (s *Store) IngestPath(ctx context.Context, req IngestRequest) (IngestResponse, error) {
	path := strings.TrimSpace(req.Path)
	if path == "" {
		return IngestResponse{}, errors.New("ingest path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return IngestResponse{}, err
	}
	run := IngestionRun{ID: newID("ing"), SourcePath: abs, Recursive: req.Recursive, Parser: "native-text", Status: "running", CreatedAt: time.Now().UTC()}
	if err := s.UpsertIngestionRun(ctx, run); err != nil {
		return IngestResponse{}, err
	}
	resp := IngestResponse{Run: run}
	finish := func(status string, runErr error) (IngestResponse, error) {
		now := time.Now().UTC()
		run.Status = status
		run.FilesSeen = len(resp.Documents) + len(resp.Skipped)
		run.DocumentsCreated = len(resp.Documents)
		run.ChunksCreated = len(resp.Chunks)
		run.CompletedAt = &now
		if runErr != nil {
			run.Error = runErr.Error()
		}
		_ = s.UpsertIngestionRun(ctx, run)
		resp.Run = run
		payload, _ := json.Marshal(map[string]any{"run_id": run.ID, "path": run.SourcePath, "recursive": run.Recursive, "status": run.Status, "documents": run.DocumentsCreated, "chunks": run.ChunksCreated, "skipped": resp.Skipped})
		_ = s.AppendEvent(ctx, Event{Kind: "document.ingested", Payload: string(payload), Source: Source{Kind: "ingest", Ref: run.ID}})
		return resp, runErr
	}

	files, err := ingestFileList(abs, req.Recursive)
	if err != nil {
		return finish("error", err)
	}
	for _, file := range files {
		if !isTextIngestFile(file) {
			resp.Skipped = append(resp.Skipped, file+" (unsupported; Docling adapter pending)")
			continue
		}
		doc, chunks, err := s.ingestTextFile(ctx, run.ID, file)
		if err != nil {
			resp.Skipped = append(resp.Skipped, file+" ("+err.Error()+")")
			continue
		}
		resp.Documents = append(resp.Documents, doc)
		resp.Chunks = append(resp.Chunks, chunks...)
	}
	status := "ok"
	if len(resp.Documents) == 0 && len(resp.Skipped) > 0 {
		status = "partial"
	}
	return finish(status, nil)
}

func ingestFileList(path string, recursive bool) ([]string, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return []string{path}, nil
	}
	var files []string
	if recursive {
		err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				name := d.Name()
				if strings.HasPrefix(name, ".") && p != path {
					return filepath.SkipDir
				}
				return nil
			}
			files = append(files, p)
			return nil
		})
	} else {
		entries, readErr := os.ReadDir(path)
		if readErr != nil {
			return nil, readErr
		}
		for _, e := range entries {
			if !e.IsDir() {
				files = append(files, filepath.Join(path, e.Name()))
			}
		}
	}
	sort.Strings(files)
	return files, err
}

func isTextIngestFile(path string) bool { return textIngestExts[strings.ToLower(filepath.Ext(path))] }

func (s *Store) ingestTextFile(ctx context.Context, runID, path string) (Document, []Chunk, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Document{}, nil, err
	}
	text := strings.TrimSpace(string(b))
	if text == "" {
		return Document{}, nil, errors.New("empty document")
	}
	sum := sha256.Sum256(b)
	hash := hex.EncodeToString(sum[:])
	doc := Document{ID: "doc_" + hash[:32], Path: path, Title: filepath.Base(path), SourceKind: "file", SourceRef: path, SHA256: hash}
	doc, err = s.UpsertDocument(ctx, doc)
	if err != nil {
		return Document{}, nil, err
	}
	if err := s.ReplaceDocumentChunks(ctx, doc.ID, splitChunks(text, 1800)); err != nil {
		return Document{}, nil, err
	}
	chunks, err := s.ListChunks(ctx, doc.ID)
	if err != nil {
		return Document{}, nil, err
	}
	return doc, chunks, nil
}

func splitChunks(text string, maxRunes int) []string {
	parts := strings.Split(text, "\n\n")
	var out []string
	var cur strings.Builder
	flush := func() {
		chunk := strings.TrimSpace(cur.String())
		if chunk != "" {
			out = append(out, chunk)
		}
		cur.Reset()
	}
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if cur.Len() > 0 && cur.Len()+len(part)+2 > maxRunes {
			flush()
		}
		if len([]rune(part)) > maxRunes {
			flush()
			runes := []rune(part)
			for len(runes) > 0 {
				n := maxRunes
				if len(runes) < n {
					n = len(runes)
				}
				out = append(out, strings.TrimSpace(string(runes[:n])))
				runes = runes[n:]
			}
			continue
		}
		if cur.Len() > 0 {
			cur.WriteString("\n\n")
		}
		cur.WriteString(part)
	}
	flush()
	if len(out) == 0 {
		out = []string{fmt.Sprintf("%s", strings.TrimSpace(text))}
	}
	return out
}
