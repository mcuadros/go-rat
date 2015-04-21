package rat

import (
	"archive/tar"
	"io"
)

type Writer struct {
	w *TrackedWriter
	t *tar.Writer
	i *Index
}

func NewWriter(w io.Writer) *Writer {
	tracked := &TrackedWriter{w, 0}
	return &Writer{
		w: tracked,
		t: tar.NewWriter(tracked),
		i: NewIndex(),
	}
}

func (w *Writer) Close() error {
	if err := w.t.Close(); err != nil {
		return err
	}

	return w.i.WriteTo(w.w)
}

func (w *Writer) Flush() error {
	return w.t.Flush()
}

func (w *Writer) Write(b []byte) (int, error) {
	return w.t.Write(b)
}

func (w *Writer) WriteHeader(hdr *tar.Header) error {
	headerPosition := w.w.position

	err := w.t.WriteHeader(hdr)
	if err != nil {
		return err
	}

	w.i.Entries[hdr.Name] = &IndexEntry{
		Name:     hdr.Name,
		Typeflag: hdr.Typeflag,
		Header:   headerPosition,
		Start:    w.w.position,
		End:      w.w.position + hdr.Size,
	}

	return err
}

type TrackedWriter struct {
	w        io.Writer
	position int64
}

func (w *TrackedWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.position += int64(n)

	return n, err
}
