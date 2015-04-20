package rat

import (
	"archive/tar"
	"io"
)

type Writer struct {
	w             *WriterWrapper
	t             *tar.Writer
	i             *Index
	currentHeader *tar.Header
	position      int64
}

func NewWriter(w io.Writer) *Writer {
	ww := &WriterWrapper{w, 0}
	return &Writer{
		w: ww,
		t: tar.NewWriter(ww),
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
		Name:   hdr.Name,
		Header: headerPosition,
		Start:  w.w.position,
		End:    w.w.position + hdr.Size,
	}

	return err
}

type WriterWrapper struct {
	w        io.Writer
	position int64
}

func (w *WriterWrapper) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.position += int64(n)

	return n, err
}
