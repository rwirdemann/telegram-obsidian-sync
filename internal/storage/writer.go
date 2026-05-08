package storage

import (
	"os"
	"path/filepath"

	"github.com/rwirdemann/telegram-obsidian-sync/internal/markdown"
)

// Writer saves Notes to a configured inbox directory.
type Writer struct {
	inboxDir string
}

// NewWriter creates a Writer targeting inboxDir.
func NewWriter(inboxDir string) *Writer {
	return &Writer{inboxDir: inboxDir}
}

// Write persists note to the inbox directory.
func (w *Writer) Write(note markdown.Note) error {
	if err := os.MkdirAll(w.inboxDir, 0755); err != nil {
		return err
	}
	path := filepath.Join(w.inboxDir, note.Filename)
	return os.WriteFile(path, []byte(note.Content), 0644)
}
