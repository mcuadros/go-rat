package rat

import (
	"archive/tar"
	"io"
)

// AddIndexToTar reads from input a standard tar file and writes on output
// a new tar file with the rat signature on it.
func AddIndexToTar(input io.Reader, output io.Writer) error {
	r := tar.NewReader(input)
	w := NewWriter(output)

	for {
		hdr, err := r.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err := w.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := io.Copy(w, r); err != nil {
			return err
		}
	}

	return w.Close()
}
