package rtar

import (
	"encoding/json"
	"io"
)

type Reader struct {
	r io.ReadSeeker
	m map[string]*Index
}

func NewReader(tarReader io.ReadSeeker, mapReader io.Reader) (*Reader, error) {
	m, err := readIndex(mapReader)
	if err != nil {
		return nil, err
	}

	return &Reader{
		r: tarReader,
		m: m,
	}, nil
}

func readIndex(r io.Reader) (map[string]*Index, error) {
	m := make(map[string]*Index, 0)

	dec := json.NewDecoder(r)
	for {
		var i Index
		if err := dec.Decode(&i); err == io.EOF {
			break
		} else if err != nil {
			return m, err
		}

		m[i.Name] = &i
	}

	return m, nil
}

func (r *Reader) ReadFile(file string) ([]byte, error) {
	i, ok := r.m[file]
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
