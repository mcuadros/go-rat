package rat

import (
	"io"
)

type Reader struct {
	r io.ReadSeeker
	i *Index
}

func NewReader(r io.ReadSeeker) (*Reader, error) {
	i := &Index{}
	if err := i.ReadFrom(r); err != nil {
		return nil, err
	}

	return &Reader{
		r: r,
		i: i,
	}, nil
}

func (r *Reader) ReadFile(file string) ([]byte, error) {
	i, ok := r.i.Entries[file]
	if !ok {
		return nil, io.EOF
	}

	content := make([]byte, i.End-i.Start)
	r.r.Seek(i.Start, 0)

	_, err := r.r.Read(content)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return content, nil
}
