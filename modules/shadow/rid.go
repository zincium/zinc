package shadow

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

// RID type
type RID [16]byte

// var s
var (
	ZeroRID RID // empty RID, all zeros
)

var rander = rand.Reader // random function

// NewRandom random
func NewRandom() (RID, error) {
	return NewRandomFromReader(rander)
}

// NewRandomFromReader returns a UUID based on bytes read from a given io.Reader.
func NewRandomFromReader(r io.Reader) (RID, error) {
	var rid RID
	_, err := io.ReadFull(r, rid[:])
	if err != nil {
		return ZeroRID, err
	}
	rid[6] = (rid[6] & 0x0f) | 0x40 // Version 4
	rid[8] = (rid[8] & 0x3f) | 0x80 // Variant is 10
	return rid, nil
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (rid RID) String() string {
	var buf [36]byte
	encodeHex(buf[:], rid)
	return string(buf[:])
}

func encodeHex(dst []byte, rid RID) {
	hex.Encode(dst, rid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], rid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], rid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], rid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], rid[10:])
}

// NewRID return RequestID
func NewRID() string {
	rid, _ := NewRandom()
	return rid.String()
}
