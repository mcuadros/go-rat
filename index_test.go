package rat

import (
	"archive/tar"
	"bytes"
	"encoding/hex"

	. "gopkg.in/check.v1"
)

var indexFixture = "524154010000000000000003000000666f6f322a0000000000000054000000000000007e0000000000000003000000626172302a0000000000000054000000000000007e000000000000004b00000000000000"
var entryFixture = "03000000666f6f302a0000000000000054000000000000007e00000000000000"

func (s *TestSuite) TestIndex_WriteTo(c *C) {
	i := index{
		map[string]*indexEntry{
			"foo": {Name: "foo", Header: 42, Start: 42 * 2, End: 42 * 3, Typeflag: tar.TypeSymlink},
			"bar": {Name: "bar", Header: 42, Start: 42 * 2, End: 42 * 3, Typeflag: tar.TypeReg},
		},
	}

	buf := bytes.NewBuffer(nil)
	err := i.WriteTo(buf)
	c.Assert(err, IsNil)

	fixture, _ := hex.DecodeString(indexFixture)
	c.Assert(buf.Bytes(), DeepEquals, fixture)
}

func (s *TestSuite) TestIndex_ReadFrom(c *C) {
	i := index{}

	fixture, _ := hex.DecodeString(indexFixture)
	buf := bytes.NewReader(fixture)
	err := i.ReadFrom(buf)
	c.Assert(err, IsNil)

	c.Assert(i.Entries, HasLen, 2)
	c.Assert(i.Entries["foo"].Name, Equals, "foo")
	c.Assert(i.Entries["bar"].Name, Equals, "bar")
}

func (s *TestSuite) TestIndexEntry_WriteTo(c *C) {
	e := indexEntry{Name: "foo", Header: 42, Start: 42 * 2, End: 42 * 3, Typeflag: tar.TypeReg}

	buf := bytes.NewBuffer(nil)
	err := e.WriteTo(buf)
	c.Assert(err, IsNil)

	fixture, _ := hex.DecodeString(entryFixture)
	c.Assert(buf.Bytes(), DeepEquals, fixture)
}

func (s *TestSuite) TestIndexEntry_WriteToInvalid(c *C) {
	e := indexEntry{}

	buf := bytes.NewBuffer(nil)
	err := e.WriteTo(buf)
	c.Assert(err, Equals, UnableToSerializeIndexEntry)
}

func (s *TestSuite) TestIndexEntry_ReadFrom(c *C) {
	e := indexEntry{}

	fixture, _ := hex.DecodeString(entryFixture)
	buf := bytes.NewBuffer(fixture)

	err := e.ReadFrom(buf)
	c.Assert(err, IsNil)
	c.Assert(e.Name, DeepEquals, "foo")
	c.Assert(e.Header, DeepEquals, int64(42))
	c.Assert(e.Start, DeepEquals, int64(42*2))
	c.Assert(e.End, DeepEquals, int64(42*3))
	c.Assert(e.Typeflag, Equals, byte(tar.TypeReg))
}

func (s *TestSuite) TestIndexEntry_ReadFromInvalid(c *C) {
	e := indexEntry{}

	buf := bytes.NewBuffer([]byte("foo"))

	err := e.ReadFrom(buf)
	c.Assert(err.Error(), Equals, "unexpected EOF")
}
