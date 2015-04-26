package bsonutil

import (
	"fmt"
	"io"

	"gopkg.in/mgo.v2/bson"
)

const MaxBSONSize = 16 * 1024 * 1024 // 16MB - maximum BSON document size

type BSONStream struct {
	err         error
	reader      io.ReadCloser
	reusableBuf []byte
}

func NewBSONStream(r io.ReadCloser) *BSONStream {
	return &BSONStream{reader: r, reusableBuf: make([]byte, MaxBSONSize)}
}

func (bs *BSONStream) Err() error {
	return bs.err
}

func (bs *BSONStream) Next(result interface{}) bool {
	hasDoc, docSize := bs.ReadNext(bs.reusableBuf)
	if !hasDoc {
		return false
	}

	if err := bson.Unmarshal(bs.reusableBuf[0:docSize], result); err != nil {
		bs.err = err
		return false
	}
	bs.err = nil
	return true
}

// ReadNext unmarshals the next BSON document into result. Returns a boolean
// indicating whether or not the operation was successful (true if no errors)
// and the size of the unmarshaled document.
func (bs *BSONStream) ReadNext(into []byte) (bool, int32) {
	// read the bson object size (a 4 byte integer)
	_, err := io.ReadAtLeast(bs.reader, into[0:4], 4)
	if err != nil {
		if err != io.EOF {
			bs.err = err
			return false, 0
		}
		// we hit EOF right away, so we're at the end of the stream.
		bs.err = nil
		return false, 0
	}

	bsonSize := int32(
		(uint32(into[0]) << 0) |
			(uint32(into[1]) << 8) |
			(uint32(into[2]) << 16) |
			(uint32(into[3]) << 24),
	)

	// Verify that the size of the BSON object we are about to read can
	// actually fit into the buffer that was provided. If not, either the BSON is
	// invalid, or the buffer passed in is too small.
	if bsonSize > int32(len(into)) {
		bs.err = fmt.Errorf("invalid BSONSize: %v bytes", bsonSize)
		return false, 0
	}
	_, err = io.ReadAtLeast(bs.reader, into[4:int(bsonSize)], int(bsonSize-4))
	if err != nil {
		if err != io.EOF {
			bs.err = err
			return false, 0
		}
		// this case means we hit EOF but read a partial document,
		// so there's a broken doc in the stream. Treat this as error.
		bs.err = fmt.Errorf("invalid bson: %v", err)
		return false, 0
	}

	bs.err = nil
	return true, bsonSize
}
