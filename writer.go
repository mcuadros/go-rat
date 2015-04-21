package rat

import (
	"archive/tar"
	"io"
)

// A Writer provides sequential writing of a rat archive. Writer has exactly the
// same interfaces as http://golang.org/pkg/archive/tar/#Writer does
type Writer struct {
	w *trackedWriter
	t *tar.Writer
	i *index
}

// NewWriter creates a new Writer writing to w.
func NewWriter(w io.Writer) *Writer {
	tracked := &trackedWriter{w, 0}
	return &Writer{
		w: tracked,
		t: tar.NewWriter(tracked),
		i: Newindex(),
	}
}

// Close closes the tar archive, flushing any unwritten data to the underlying
// writer and writes the rat signature at the end of the writer.
func (w *Writer) Close() error {
	if err := w.t.Close(); err != nil {
		return err
	}

	return w.i.WriteTo(w.w)
}

// Flush finishes writing the current file (optional).
func (w *Writer) Flush() error {
	return w.t.Flush()
}

// Write writes to the current entry in the tar archive.
func (w *Writer) Write(b []byte) (int, error) {
	return w.t.Write(b)
}

// WriteHeader writes hdr and prepares to accept the file's contents.
func (w *Writer) WriteHeader(hdr *tar.Header) error {
	headerPosition := w.w.position

	err := w.t.WriteHeader(hdr)
	if err != nil {
		return err
	}

	w.i.Entries[hdr.Name] = &indexEntry{
		Name:     hdr.Name,
		Typeflag: hdr.Typeflag,
		Header:   headerPosition,
		Start:    w.w.position,
		End:      w.w.position + hdr.Size,
	}

	return err
}

type trackedWriter struct {
	w        io.Writer
	position int64
}

func (w *trackedWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.position += int64(n)

	return n, err
}
