package rat

import (
	"bytes"
	"os"

	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestAddIndexIndexToTar(c *C) {
	f, err := os.Open("fixtures/bsd.tar")
	c.Assert(err, IsNil)

	buf := bytes.NewBuffer(nil)
	err = AddIndexIndexToTar(f, buf)
	c.Assert(err, IsNil)

	r, err := NewReader(bytes.NewReader(buf.Bytes()))
	c.Assert(err, IsNil)

	content, err := r.ReadFile("composer.json")
	c.Assert(err, IsNil)
	c.Assert(content, HasLen, 4011)
}
