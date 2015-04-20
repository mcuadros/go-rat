package rat

import (
	"archive/tar"
	"bytes"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) TestWriter(c *C) {
	f := bytes.NewBuffer(nil)
	i := bytes.NewBuffer(nil)
	w := NewWriter(f, i)

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

	r, err := NewReader(bytes.NewReader(f.Bytes()), i)
	c.Assert(err, IsNil)

	d, err := r.ReadFile("foo.txt")
	c.Assert(err, IsNil)
	c.Assert(string(d), Equals, "foo")

	d, err = r.ReadFile("bar.txt")
	c.Assert(err, IsNil)
	c.Assert(string(d), Equals, "bar")
}
