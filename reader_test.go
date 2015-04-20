package rat

import (
	"archive/tar"
	"bytes"

	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestNewReader(c *C) {
	f := bytes.NewBuffer(nil)

	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})
	w.Write([]byte("foo"))
	w.Close()

	_, err := NewReader(bytes.NewReader(f.Bytes()))
	c.Assert(err, IsNil)
}

func (s *TestSuite) TestNewReader_NotRAT(c *C) {
	f := bytes.NewBuffer(nil)

	w := tar.NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})
	w.Write([]byte("foo"))
	w.Close()

	_, err := NewReader(bytes.NewReader(f.Bytes()))
	c.Assert(err, Equals, UnsuportedIndex)
}

func (s *TestSuite) TestReader_ReadFile(c *C) {
	f := bytes.NewBuffer(nil)

	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})
	w.Write([]byte("foo"))
	w.Close()

	r, err := NewReader(bytes.NewReader(f.Bytes()))
	c.Assert(err, IsNil)

	content, err := r.ReadFile("foo.txt")
	c.Assert(err, IsNil)
	c.Assert(string(content), Equals, "foo")
}

func (s *TestSuite) TestReader_ReadFileNotFound(c *C) {
	f := bytes.NewBuffer(nil)

	w := NewWriter(f)
	w.Close()

	r, err := NewReader(bytes.NewReader(f.Bytes()))
	c.Assert(err, IsNil)

	_, err = r.ReadFile("foo.txt")
	c.Assert(err, Equals, FileNotFound)
}

func (s *TestSuite) TestReader_ReadFileInvalidIndex(c *C) {
	f := bytes.NewBuffer(nil)

	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})
	w.Write([]byte("foo"))
	w.Close()

	r, err := NewReader(bytes.NewReader(f.Bytes()))
	r.i.Entries["foo.txt"].Start -= 1000
	r.i.Entries["foo.txt"].End -= 1000
	c.Assert(err, IsNil)

	_, err = r.ReadFile("foo.txt")
	c.Assert(err, Not(IsNil))
}
