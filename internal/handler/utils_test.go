package handler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadStrmContent(t *testing.T) {
	dir := t.TempDir()

	t.Run("non-strm path returned unchanged", func(t *testing.T) {
		input := "/data/media/movie.mkv"
		if got := readStrmContent(input); got != input {
			t.Errorf("expected %q, got %q", input, got)
		}
	})

	t.Run("strm file content returned trimmed", func(t *testing.T) {
		p := filepath.Join(dir, "movie.strm")
		if err := os.WriteFile(p, []byte("  /115/movie.mkv\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if got := readStrmContent(p); got != "/115/movie.mkv" {
			t.Errorf("expected %q, got %q", "/115/movie.mkv", got)
		}
	})

	t.Run("uppercase .STRM extension handled", func(t *testing.T) {
		p := filepath.Join(dir, "movie.STRM")
		if err := os.WriteFile(p, []byte("https://example.com/video.mp4"), 0o600); err != nil {
			t.Fatal(err)
		}
		if got := readStrmContent(p); got != "https://example.com/video.mp4" {
			t.Errorf("expected %q, got %q", "https://example.com/video.mp4", got)
		}
	})

	t.Run("missing strm file returns path unchanged", func(t *testing.T) {
		p := filepath.Join(dir, "missing.strm")
		if got := readStrmContent(p); got != p {
			t.Errorf("expected %q, got %q", p, got)
		}
	})
}
