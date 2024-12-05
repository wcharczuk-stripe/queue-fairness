package sim

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync/atomic"
)

var uuidBase uint64 = 1

func NewUUID() UUID {
	newUID := atomic.AddUint64(&uuidBase, 1)
	var uuid UUID
	binary.PutUvarint(uuid[:], newUID)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // set version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // set variant 2
	return uuid
}

// UUID represents a unique identifier conforming to the RFC 4122 standard.
// UUIDs are a fixed 128bit (16 byte) binary blob.
type UUID [16]byte

// Equal returns if a uuid is equal to another uuid.
func (uuid UUID) Equal(other UUID) bool {
	return bytes.Equal(uuid[0:], other[0:])
}

// Compare returns a comparison between the two uuids.
func (uuid UUID) Compare(other UUID) int {
	return bytes.Compare(uuid[0:], other[0:])
}

// String returns the uuid as a hex string.
func (uuid UUID) String() string {
	return hex.EncodeToString([]byte(uuid[:]))
}

// ShortString returns the first 8 bytes of the uuid as a hex string.
func (uuid UUID) ShortString() string {
	return hex.EncodeToString([]byte(uuid[:8]))
}

// Version returns the version byte of a uuid.
func (uuid UUID) Version() byte {
	return uuid[6] >> 4
}

// Format allows for conditional expansion in printf statements
// based on the token and flags used.
func (uuid UUID) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			b := []byte(uuid[:])
			fmt.Fprintf(s,
				"%08x-%04x-%04x-%04x-%012x",
				b[:4], b[4:6], b[6:8], b[8:10], b[10:],
			)
			return
		}
		fmt.Fprint(s, hex.EncodeToString([]byte(uuid[:])))
	case 's':
		fmt.Fprint(s, hex.EncodeToString([]byte(uuid[:])))
	case 'q':
		fmt.Fprintf(s, "%b", uuid.Version())
	}
}

// IsZero returns if the uuid is unset.
func (uuid UUID) IsZero() bool {
	return uuid == [16]byte{}
}
