package rat

import (
	"archive/tar"
	"errors"
	"io"
)

type Reader struct {
	r io.ReadSeeker
	i *Index
}

var (
	FileNotFound = errors.New("File not found")
	NotRegFile   = errors.New("This is not a regular file")
)

func NewReader(r io.ReadSeeker) (*Reader, error) {
	i := &Index{}
	if err := i.ReadFrom(r); err != nil {
		return nil, err
	}

	return &Reader{r: r, i: i}, nil
}

func (r *Reader) GetNames(onlyRegFiles bool) []string {
	result := make([]string, 0)

	for name, i := range r.i.Entries {
		if onlyRegFiles && i.Typeflag != tar.TypeReg && i.Typeflag != tar.TypeRegA {
			continue
		}

		result = append(result, name)
	}

	return result
}

func (r *Reader) ReadFile(file string) ([]byte, error) {
	i, ok := r.i.Entries[file]
	if !ok {
		return nil, FileNotFound
	}

	if i.Typeflag != tar.TypeReg && i.Typeflag != tar.TypeRegA {
		return nil, NotRegFile
	}

	if _, err := r.r.Seek(i.Start, 0); err != nil {
		return nil, err
	}

	content := make([]byte, i.End-i.Start)
	if _, err := r.r.Read(content); err != nil && err != io.EOF {
		return nil, err
	}

	return content, nil
}
