package rat

import (
	"archive/tar"
	"io"
)

type Writer struct {
	w        *WriterWrapper
	t        *tar.Writer
	i        *IndexWriter
	index    *Index
	position int64
}

func NewWriter(tarWriter, mapWriter io.Writer) *Writer {
	w := &WriterWrapper{tarWriter, 0}
	return &Writer{
		w: w,
		t: tar.NewWriter(w),
		i: &IndexWriter{mapWriter},
	}
}

func (w *Writer) Close() error {
	return w.t.Close()
}

func (w *Writer) Flush() error {
	return w.t.Flush()
}

func (w *Writer) Write(b []byte) (int, error) {
	n, err := w.t.Write(b)
	w.index.End = w.w.position
	w.i.Write(w.index)
	w.index = nil

	return n, err
}

func (w *Writer) WriteHeader(hdr *tar.Header) error {
	err := w.t.WriteHeader(hdr)
	w.index = &Index{
		Name:  hdr.Name,
		Start: w.w.position,
	}

	return err
}

type Index struct {
	Name       string
	Start, End int64
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
