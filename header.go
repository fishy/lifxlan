package lifxlan

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func init() {
	// Seed pseudo-random number generator used by RandomSource().
	now := time.Now()
	rand.Seed(now.Unix() + int64(now.Nanosecond()))
}

// TaggedHeader is the 16-bit header including:
// - origin: 2 bits, must be 0
// - tagged: 1 bit
// - addressable: 1 bit, must be 1
// - protocal: 12 bits, must be 1024
type TaggedHeader uint16

// Tagged and non-tagged versions of TaggedHeader.
const (
	NotTagged TaggedHeader = 1<<12 + 1024
	Tagged                 = 1<<13 + NotTagged
)

// AckResFlag could include:
// - ack_required: if set all sent messages will expect an ack response.
// - res_required: if set all sent messages will expect a response.
type AckResFlag uint8

// AckResFlag values.
const (
	FlagResRequired AckResFlag = 1 << iota
	FlagAckRequired
)

// Header offsets.
const (
	SizeOffset      = 0                   // 2 bytes
	TaggedOffset    = SizeOffset + 2      // 2 bytes
	SourceOffset    = TaggedOffset + 2    // 4 bytes
	TargetOffset    = SourceOffset + 4    // 8 bytes
	reserved1Offset = TargetOffset + 8    // 6 bytes
	FlagsOffset     = reserved1Offset + 6 // 1 byte
	SequenceOffset  = FlagsOffset + 1     // 1 byte
	reserved2Offset = SequenceOffset + 1  // 8 bytes
	TypeOffset      = reserved2Offset + 8 // 2 bytes
	reserved3Offset = TypeOffset + 2      // 2 bytes
	PayloadOffset   = reserved3Offset + 2
)

// HeaderLength is the length of the header
const HeaderLength uint16 = PayloadOffset

// ResponseReadBufferSize is the recommended buffer size to read UDP responses.
// It's big enough for all the payloads.
const ResponseReadBufferSize = 4096

// UintToBytes encodes v into a byte array with appropriate size.
//
// v must be an uintN type (excluding uint), or UintToBytes panics.
func UintToBytes(v interface{}) []byte {
	u := reflect.ValueOf(v).Uint()
	size := binary.Size(v)
	buf := make([]byte, size)
	switch size {
	case 1:
		buf[0] = byte(u)
	case 2:
		binary.LittleEndian.PutUint16(buf, uint16(u))
	case 4:
		binary.LittleEndian.PutUint32(buf, uint32(u))
	case 8:
		binary.LittleEndian.PutUint64(buf, uint64(u))
	}
	return buf
}

// GenerateMessage generates the message to send.
func GenerateMessage(
	tagged TaggedHeader,
	source uint32,
	target Target,
	flags AckResFlag,
	sequence uint8,
	message MessageType,
	payload []byte,
) []byte {
	var size = HeaderLength + uint16(len(payload))
	buf := make([]byte, size)
	var data = map[int][]byte{
		SizeOffset:     UintToBytes(size),     // size
		TaggedOffset:   UintToBytes(tagged),   // origin, tagged, addressable, protocol
		SourceOffset:   UintToBytes(source),   // source
		TargetOffset:   UintToBytes(target),   // target
		FlagsOffset:    UintToBytes(flags),    // reserved, ack_required, res_required
		SequenceOffset: UintToBytes(sequence), // sequence
		TypeOffset:     UintToBytes(message),  // type
		PayloadOffset:  payload,
	}
	for offset, v := range data {
		copy(buf[offset:], v)
	}
	return buf
}

// Response is the parsed response from a lifxlan device.
type Response struct {
	Message  MessageType
	Source   uint32
	Target   Target
	Sequence uint8
	Payload  []byte
}

// ParseResponse parses the response received from a lifxlan device.
func ParseResponse(msg []byte) (*Response, error) {
	if len(msg) < int(HeaderLength) {
		return nil, fmt.Errorf(
			"lifxlan.ParseResponse: response size not enough: %d < %d",
			len(msg),
			HeaderLength,
		)
	}
	size := binary.LittleEndian.Uint16(msg[SizeOffset:])
	if len(msg) != int(size) {
		return nil, fmt.Errorf(
			"lifxlan.ParseResponse: response size mismatch: %d != %d",
			len(msg),
			size,
		)
	}
	return &Response{
		Message:  MessageType(binary.LittleEndian.Uint16(msg[TypeOffset:])),
		Source:   binary.LittleEndian.Uint32(msg[SourceOffset:]),
		Target:   Target(binary.LittleEndian.Uint64(msg[TargetOffset:])),
		Sequence: uint8(msg[SequenceOffset]),
		Payload:  msg[PayloadOffset:],
	}, nil
}

var maxSource int64 = math.MaxUint32

// RandomSource generates a random number to be used as source.
// It's guaranteed to be non-zero when err is nil.
func RandomSource() uint32 {
	return uint32(rand.Int63n(maxSource) + 1)
}
