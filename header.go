package lifxlan

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

func init() {
	// Seed pseudo-random number generator used by RandomSource().
	now := time.Now()
	rand.Seed(now.Unix() + int64(now.Nanosecond()))
}

// TaggedHeader is the 16-bit header including:
//
// - origin: 2 bits, must be 0
//
// - tagged: 1 bit
//
// - addressable: 1 bit, must be 1
//
// - protocol: 12 bits, must be 1024
type TaggedHeader uint16

// Tagged and non-tagged versions of TaggedHeader.
const (
	NotTagged TaggedHeader = 1<<12 + 1024
	Tagged                 = 1<<13 + NotTagged
)

// AckResFlag is the 8-bit header that could include:
//
// - ack_required: if set all sent messages will expect an ack response.
//
// - res_required: if set all sent messages will expect a response.
type AckResFlag uint8

// AckResFlag values.
const (
	FlagResRequired AckResFlag = 1 << iota
	FlagAckRequired
)

// RawHeader defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/header-description
type RawHeader struct {
	Size     uint16
	Tagged   TaggedHeader
	Source   uint32
	Target   Target
	_        [6]uint8 // reserved
	Flags    AckResFlag
	Sequence uint8
	_        uint64 // reserved
	Type     MessageType
	_        uint16 // reserved
}

// HeaderLength is the length of the header
const HeaderLength = 36

// ResponseReadBufferSize is the recommended buffer size to read UDP responses.
// It's big enough for all the payloads.
const ResponseReadBufferSize = 4096

// GenerateMessage generates the message to send.
func GenerateMessage(
	tagged TaggedHeader,
	source uint32,
	target Target,
	flags AckResFlag,
	sequence uint8,
	message MessageType,
	payload []byte,
) ([]byte, error) {
	var size = HeaderLength + uint16(len(payload))
	buf := new(bytes.Buffer)
	data := &RawHeader{
		Size:     size,
		Tagged:   tagged,
		Source:   source,
		Target:   target,
		Flags:    flags,
		Sequence: sequence,
		Type:     message,
	}
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	if _, err := buf.Write(payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

	var d RawHeader
	r := bytes.NewReader(msg)
	if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
		return nil, err
	}
	if len(msg) != int(d.Size) {
		return nil, fmt.Errorf(
			"lifxlan.ParseResponse: response size mismatch: %d != %d",
			len(msg),
			d.Size,
		)
	}

	payload, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Response{
		Message:  d.Type,
		Source:   d.Source,
		Target:   d.Target,
		Sequence: d.Sequence,
		Payload:  payload,
	}, nil
}

var maxSource int64 = math.MaxUint32

// RandomSource generates a random number to be used as source.
// It's guaranteed to be non-zero.
func RandomSource() uint32 {
	return uint32(rand.Int63n(maxSource) + 1)
}
