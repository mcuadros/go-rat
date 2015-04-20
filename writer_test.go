package rat

import (
	"archive/tar"
	"bytes"

	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestWriter_WriteHeader(c *C) {
	f := bytes.NewBuffer(nil)
	w := NewWriter(f)
	c.Assert(w.i.Entries, HasLen, 0)

	err := w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})
	c.Assert(err, IsNil)

	c.Assert(w.i.Entries, HasLen, 1)
	e := w.i.Entries["foo.txt"]
	c.Assert(e.Name, Equals, "foo.txt")
	c.Assert(e.Header, Equals, int64(0))
	c.Assert(e.Start, Equals, int64(512))
	c.Assert(e.End, Equals, int64(515))
}

func (s *TestSuite) TestWriter_Write(c *C) {
	f := bytes.NewBuffer(nil)
	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})

	n, err := w.Write([]byte("foo"))
	c.Assert(err, IsNil)
	c.Assert(n, Equals, 3)
}

func (s *TestSuite) TestWriter_Flush(c *C) {
	f := bytes.NewBuffer(nil)
	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})

	err := w.Flush()
	c.Assert(err, Not(IsNil))
}

func (s *TestSuite) TestWriter_Close(c *C) {
	f := bytes.NewBuffer(nil)
	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})
	w.Write([]byte("foo"))

	c.Assert(f.Len(), Equals, 515)
	err := w.Close()
	c.Assert(err, IsNil)
	c.Assert(f.Len(), Equals, 2102)
}

func (s *TestSuite) TestWriter_CloseFailed(c *C) {
	f := bytes.NewBuffer(nil)
	w := NewWriter(f)
	w.WriteHeader(&tar.Header{Name: "foo.txt", Size: 3})

	err := w.Close()
	c.Assert(err, Not(IsNil))
}
