package rat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const IndexVersion int64 = 1

var (
	IndexSignature              = []byte{'R', 'A', 'T'}
	UnsuportedIndexVersion      = errors.New("Unsuported Index signature")
	WrongIndexSignature         = errors.New("Wrong Index signature")
	UnableToSerializeIndexEntry = errors.New("Unable to serialize: invalid content")
)

type Index struct {
	Entries map[string]*IndexEntry
}

func NewIndex() *Index {
	return &Index{make(map[string]*IndexEntry, 0)}
}

// IndexEntry byte representation on LittleEndian have the following format:
// - 3-byte index signature
// - x-byte index entries
// - 8-byte length
func (i *Index) WriteTo(w io.Writer) error {
	tail := bytes.NewBuffer(IndexSignature)
	if err := binary.Write(tail, binary.LittleEndian, IndexVersion); err != nil {
		return err
	}

	for _, e := range i.Entries {
		if err := e.WriteTo(tail); err != nil {
			return err
		}
	}

	length, err := io.Copy(w, tail)
	if err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, int64(length)); err != nil {
		return err
	}

	return nil
}

const tailSizeLength = 8 //int64

func (i *Index) ReadFrom(r io.ReadSeeker) error {
	if _, err := r.Seek(-tailSizeLength, 2); err != nil {
		return err
	}

	var tailLen int64
	if err := binary.Read(r, binary.LittleEndian, &tailLen); err != nil {
		return err
	}

	if _, err := r.Seek(-tailSizeLength-tailLen, 2); err != nil {
		return err
	}

	sig := make([]byte, 3)
	if _, err := r.Read(sig); err != nil {
		return err
	}

	if !bytes.Equal(sig, IndexSignature) {
		return WrongIndexSignature
	}

	var version int64
	if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
		return err
	}

	if version != IndexVersion {
		return UnsuportedIndexVersion
	}

	i.Entries = make(map[string]*IndexEntry, 0)
	for {
		e := &IndexEntry{}
		if err := e.ReadFrom(r); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		i.Entries[e.Name] = e
	}

	return nil
}

type IndexEntry struct {
	Name       string
	Header     int64
	Start, End int64
}

// IndexEntry byte representation on LittleEndian have the following format:
// - 4-byte length of the filename
// - x-byte filename
// - 8-byte header
// - 8-byte start
// - 8-byte end
func (i *IndexEntry) WriteTo(w io.Writer) error {
	if i.Name == "" || i.Start == 0 || i.End == 0 {
		return UnableToSerializeIndexEntry
	}

	name := []byte(i.Name)
	if err := binary.Write(w, binary.LittleEndian, int32(len(name))); err != nil {
		return err
	}

	if _, err := w.Write(name); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, i.Header); err != nil {
		return err
	}

	//TODO: not allow 0 in start or end
	if err := binary.Write(w, binary.LittleEndian, i.Start); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, i.End); err != nil {
		return err
	}

	return nil
}

func (i *IndexEntry) ReadFrom(r io.Reader) error {
	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return err
	}

	filename := make([]byte, length)
	if _, err := r.Read(filename); err != nil {
		return err
	}

	i.Name = string(filename)

	err := binary.Read(r, binary.LittleEndian, &i.Header)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &i.Start)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &i.End)
	if err != nil {
		return err
	}

	return nil
}
