package rat

import (
	"archive/tar"
	"bytes"
	"os"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestWriter(c *C) {
	f := bytes.NewBuffer(nil)
	w := NewWriter(f)

	h := &tar.Header{
		Name:    "foo.txt",
		Mode:    0640,
		Size:    3,
		ModTime: time.Unix(1244428340, 0),
	}

	c.Assert(w.WriteHeader(h), IsNil)
	n, err := w.Write([]byte("foo"))
	c.Assert(n, Equals, 3)
	c.Assert(err, IsNil)

	c.Assert(w.WriteHeader(&tar.Header{
		Name:    "bar.txt",
		Mode:    0640,
		Size:    3,
		ModTime: time.Unix(1244428340, 0),
	}), IsNil)

	n, err = w.Write([]byte("bar"))
	c.Assert(n, Equals, 3)
	c.Assert(err, IsNil)

	c.Assert(w.Close(), IsNil)

	r, err := NewReader(bytes.NewReader(f.Bytes()))
	c.Assert(err, IsNil)

	d, err := r.ReadFile("foo.txt")
	c.Assert(err, IsNil)
	c.Assert(string(d), Equals, "foo")

	d, err = r.ReadFile("bar.txt")
	c.Assert(err, IsNil)
	c.Assert(string(d), Equals, "bar")
}

func (s *TestSuite) TestWriterRealFile(c *C) {
	f, err := os.Create("/tmp/foo.tar")
	c.Assert(err, IsNil)

	w := NewWriter(f)

	h := &tar.Header{
		Name:    "foo.txt",
		Mode:    0640,
		Size:    3,
		ModTime: time.Unix(1244428340, 0),
	}

	c.Assert(w.WriteHeader(h), IsNil)
	n, err := w.Write([]byte("foo"))
	c.Assert(n, Equals, 3)
	c.Assert(err, IsNil)

	c.Assert(w.WriteHeader(&tar.Header{
		Name:    "bar.txt",
		Mode:    0640,
		Size:    3,
		ModTime: time.Unix(1244428340, 0),
	}), IsNil)

	n, err = w.Write([]byte("bar"))
	c.Assert(n, Equals, 3)
	c.Assert(err, IsNil)

	c.Assert(w.Close(), IsNil)

	f, err = os.Open("/tmp/foo.tar")
	c.Assert(err, IsNil)

	r, err := NewReader(f)
	c.Assert(err, IsNil)

	d, err := r.ReadFile("foo.txt")
	c.Assert(err, IsNil)
	c.Assert(string(d), Equals, "foo")

	d, err = r.ReadFile("bar.txt")
	c.Assert(err, IsNil)
	c.Assert(string(d), Equals, "bar")
}
