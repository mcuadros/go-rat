package rat

import (
	"encoding/binary"
	"io"
)

type IndexWriter struct {
	w io.Writer
}

// Index byte representation on LittleEndian have the following format:
// - 4-byte length of the filename
// - 4-byte filename
// - 8-byte start
// - 8-byte end
func (w *IndexWriter) Write(i *Index) error {
	name := []byte(i.Name)
	if err := binary.Write(w.w, binary.LittleEndian, int32(len(name))); err != nil {
		return err
	}

	if _, err := w.w.Write(name); err != nil {
		return err
	}

	if err := binary.Write(w.w, binary.LittleEndian, i.Start); err != nil {
		return err
	}

	if err := binary.Write(w.w, binary.LittleEndian, i.End); err != nil {
		return err
	}

	return nil
}

type IndexReader struct {
	r io.Reader
}

func (r *IndexReader) Read() (*Index, error) {
	var length int32
	if err := binary.Read(r.r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	filename := make([]byte, length)
	if _, err := r.r.Read(filename); err != nil {
		return nil, err
	}

	i := &Index{Name: string(filename)}

	if err := binary.Read(r.r, binary.LittleEndian, &i.Start); err != nil {
		return nil, err
	}

	if err := binary.Read(r.r, binary.LittleEndian, &i.End); err != nil {
		return nil, err
	}

	return i, nil
}
