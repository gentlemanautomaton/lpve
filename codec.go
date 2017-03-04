package lpve

import "io"

// Codec is a length-ordered value encoder/decoder for a particular hash size.
type Codec struct {
	// TODO: Inject filter and hash generation functions?
	Length int // Hash length in bytes, between 1 and 8
}

var (
	// Hash128 is an lpve codec for 128 bit hash values.
	Hash128 = Codec{16}

	// Hash256 is an lpve codec for 256 bit hash values.
	Hash256 = Codec{32}

	// Hash512 is an lpve codec for 512 bit hash values.
	Hash512 = Codec{64}
)

// Value is a variable-length encoding for binary octet strings.
type Value []byte

// Content represents a fixed width space for variable-length content.
type Content [8]byte

// Constants
const (
	MaxLength = 0x7FFFFFFFFFFFFFF // Maximum length of a lpve-encodable octet string.
)

// Encoding type identification
const (
	TypeMask = 0xC0 // 1100 0000

	// TypeNil designates a zero length octet string.
	TypeNil = 0x00 // 0000 0000

	// TypeInlineByte designates an octet string of length 1 with an unsigned
	// value not exceeding 63.
	TypeInlineByte = 0x40 // 0100 0000

	// TypeInlineMultibyte designates an octet string of length 1-64 with an
	// unsigned value not exceeding 63.
	TypeInlineMultibyte = 0x80 // 1000 0000

	// TypeReference designates an octet string of length greater than 64
	// that is represented by a reference, which is a sakura .
	TypeReference = 0xC0 // 1100 0000
)

// Length masks
const (
	InlineMask             = 0x3F // 0011 1111
	ReferencePrefixMask    = 0x38 // 0011 1000
	ReferenceCarryoverMask = 0x07 // 0000 0111
)

// Reference length offsets
const (
	offset1 = 0x000000000000008 + 64 + 1
	offset2 = 0x000000000000808 + 64 + 1
	offset3 = 0x000000000080808 + 64 + 1
	offset4 = 0x000000008080808 + 64 + 1
	offset5 = 0x000000808080808 + 64 + 1
	offset6 = 0x000080808080808 + 64 + 1
	offset7 = 0x008080808080808 + 64 + 1
)

// ParseValue parses the given byte slice
func ParseValue(b []byte) Value {
	return nil
}

// Len returns the uncompressed total length of the data represented by value,
// in bytes.
func (v Value) Len() uint64 {
	switch v[0] & TypeMask {
	case TypeNil:
		return 0
	case TypeInlineByte:
		return 1
	case TypeInlineMultibyte:
		return uint64(v[0] & InlineMask)
	case TypeReference:
		d := uint64((v[0] & ReferencePrefixMask) >> 3)
		j := uint64((v[0] & ReferenceCarryoverMask))
		length := 64 + (1 << (8 * d)) + d
		length |= j << (8 * d)
		for i := uint64(0); i < d; i++ {
			length |= uint64(v[i]) << 8 * i
		}
		return length
	}
	return 0
}

func encode(w io.Writer) (length, value uint64, err error) {
	return 0, 0, nil
}

func decode(r io.Reader) (length, content Content, err error) {
	//var buf [8]byte
	//var header byte
	return
}

// DecodeSlice will decode the given byte slice and return the encoded length
// and content.
func (c Codec) DecodeSlice(b []byte) (length uint64, content Content) {
	header := b[0]

	// Inline types

	switch header & TypeMask {
	case TypeNil:
		return
	case TypeInlineByte:
		length = 1
		content[0] = header & InlineMask
		return
	case TypeInlineMultibyte:
		length = uint64(header&InlineMask + 1)
		switch length {
		case 1:
			content[0] = b[1]
		case 2:
			content[0], content[1] = b[1], b[2]
		case 3:
			content[0], content[1], content[2] = b[1], b[2], b[3]
		case 4:
			content[0], content[1], content[2], content[3] = b[1], b[2], b[3], b[4]
		case 5:
			content[0], content[1], content[2], content[3], content[4] = b[1], b[2], b[3], b[4], b[5]
		case 6:
			content[0], content[1], content[2], content[3], content[4], content[5] = b[1], b[2], b[3], b[4], b[5], b[6]
		case 7:
			content[0], content[1], content[2], content[3], content[4], content[5], content[6] = b[1], b[2], b[3], b[4], b[5], b[6], b[7]
		}
		return
	}

	// Reference length
	count := (header & ReferencePrefixMask) >> 3     // 00XXX000
	length = uint64(header & ReferenceCarryoverMask) // 00000XXX

	switch count {
	case 0:
	case 1:
		length = length<<8 | uint64(b[1]) + offset1
	case 2:
		length = length<<16 | uint64(b[1])<<8 | uint64(b[2]) + offset2
	case 3:
		length = length<<24 | uint64(b[1])<<16 | uint64(b[2])<<8 | uint64(b[3]) + offset3
	case 4:
		length = length<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4]) + offset4
	case 5:
		length = length<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]) + offset5
	case 6:
		length = length<<48 | uint64(b[1])<<40 | uint64(b[2])<<32 | uint64(b[3])<<24 | uint64(b[4])<<16 | uint64(b[5])<<8 | uint64(b[6]) + offset6
	case 7:
		length = length<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 | uint64(b[4])<<24 | uint64(b[5])<<8 | uint64(b[6])<<8 | uint64(b[7]) + offset7
	}

	// Reference content
	b = b[1+count:]
	if len(b) < c.Length {
		panic("invalid length prefix encoding")
	}

	for i := 0; i < c.Length; i++ {
		content[i] = b[i]
	}

	return
}

// Extract stores the value or reference encoded  returns the uncompressed total length of the data represented by value,
// in bytes.
func (v Value) Extract(b *[72]byte) uint64 {
	switch v[0] & TypeMask {
	case TypeNil:
		return 0
	case TypeInlineByte:
		return uint64(v[0] & InlineMask)
	case TypeInlineMultibyte:
		return 0
	}
	return 0
}

// Inline returns true if the value is stored within the value itself.
func (v Value) Inline() bool {
	return false
}

// Bytes returns the bytes contained in the value if the data is inlined.
func (v Value) Bytes() (b []byte, ok bool) {
	return nil, false
}
