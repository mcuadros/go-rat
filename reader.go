package rat

import (
	"archive/tar"
	"errors"
	"io"
)

// A Reader provides random access to the contents of a rat archive. You can use
// archive/tar.Reader with any rat file for sequencial access.
type Reader struct {
	r io.ReadSeeker
	i *index
}

var (
	FileNotFound = errors.New("File not found")
	NotRegFile   = errors.New("This is not a regular file")
)

// NewReader creates a new Reader reading from r.
func NewReader(r io.ReadSeeker) (*Reader, error) {
	i := &index{}
	if err := i.ReadFrom(r); err != nil {
		return nil, err
	}

	return &Reader{r: r, i: i}, nil
}

// GetNames returns all the entries from the rat signautre, you can filter only
// the regular files using the onlyRegFiles arg
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

// ReadFile returns the content of a file, if the entry not is a regular file
// the error NotRegFile is returned, if the files not exists returns FileNotFound
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
