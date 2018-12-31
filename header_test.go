package lifxlan_test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"testing"
	"testing/quick"

	"github.com/fishy/lifxlan"
)

func TestParseResponseError(t *testing.T) {
	makeMsg := func(size uint16) []byte {
		buf := make([]byte, int(size))
		binary.LittleEndian.PutUint16(buf, size)
		return buf
	}

	t.Run(
		"SizeNotEnough",
		func(t *testing.T) {
			t.Run(
				"Empty",
				func(t *testing.T) {
					var buf []byte
					_, err := lifxlan.ParseResponse(buf)
					if err == nil {
						t.Errorf("Expected size not enough error for msg % x", buf)
					}
				},
			)

			t.Run(
				"NonEmpty",
				func(t *testing.T) {
					size := lifxlan.HeaderLength - 1
					buf := makeMsg(size)
					_, err := lifxlan.ParseResponse(buf)
					if err == nil {
						t.Errorf("Expected size not enough error for msg % x", buf)
					}
				},
			)
		},
	)

	t.Run(
		"SizeMismatch",
		func(t *testing.T) {
			size := lifxlan.HeaderLength + 10
			buf := makeMsg(size)[:lifxlan.HeaderLength+1]
			_, err := lifxlan.ParseResponse(buf)
			if err == nil {
				t.Errorf("Expected size mismatch error for msg % x", buf)
			}
		},
	)
}

func TestHeader(t *testing.T) {
	var sequence uint8

	t.Run(
		"Discover",
		func(t *testing.T) {
			sequence++

			msg := lifxlan.GenerateMessage(
				lifxlan.Tagged,
				0, // source
				lifxlan.AllDevices,
				0, // flags
				sequence,
				lifxlan.GetService,
				nil, // payload
			)
			resp, err := lifxlan.ParseResponse(msg)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Message != lifxlan.GetService {
				t.Errorf(
					"resp.Message expected %v, got %v",
					lifxlan.GetService,
					resp.Message,
				)
			}
			if resp.Source != 0 {
				t.Errorf("resp.Source expected 0, got %v", resp.Source)
			}
			if resp.Target != lifxlan.AllDevices {
				t.Errorf(
					"resp.Target expected %v, got %v",
					lifxlan.AllDevices,
					resp.Target,
				)
			}
			if resp.Sequence != sequence {
				t.Errorf("resp.Sequence expected %v, got %v", sequence, resp.Sequence)
			}
			if len(resp.Payload) != 0 {
				t.Errorf("resp.Payload expected empty, got % x", resp.Payload)
			}
		},
	)

	t.Run(
		"WithPayload",
		func(t *testing.T) {
			sequence++
			source := lifxlan.RandomSource()
			target := lifxlan.Target(1234)
			payload := make([]byte, 10)
			_, err := rand.Reader.Read(payload)
			if err != nil {
				t.Fatal(err)
			}
			msgType := lifxlan.MessageType(4321)

			msg := lifxlan.GenerateMessage(
				lifxlan.NotTagged,
				source,
				target,
				0, // flags
				sequence,
				msgType,
				payload,
			)
			resp, err := lifxlan.ParseResponse(msg)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Message != msgType {
				t.Errorf(
					"resp.Message expected %v, got %v",
					msgType,
					resp.Message,
				)
			}
			if resp.Source != source {
				t.Errorf("resp.Source expected %v, got %v", source, resp.Source)
			}
			if resp.Target != target {
				t.Errorf(
					"resp.Target expected %v, got %v",
					target,
					resp.Target,
				)
			}
			if resp.Sequence != sequence {
				t.Errorf("resp.Sequence expected %v, got %v", sequence, resp.Sequence)
			}
			if !bytes.Equal(payload, resp.Payload) {
				t.Errorf("resp.Payload expected % x, got % x", payload, resp.Payload)
			}
		},
	)
}

func TestRandomSource(t *testing.T) {
	n := 0
	f := func() bool {
		n++
		v := lifxlan.RandomSource()
		if v == 0 {
			t.Log("RandomSource returned 0")
			return false
		}
		t.Logf("RandomSource returned %v", v)
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
	t.Logf("quick did %d checks", n)
}

func TestUintToBytes(t *testing.T) {
	t.Run(
		"panic",
		func(t *testing.T) {
			defer func() {
				if err := recover(); err == nil {
					t.Error("UintToBytes(non-uint) didn't panic")
				}
			}()
			lifxlan.UintToBytes(int(3))
		},
	)

	var value interface{}
	var expect, actual []byte

	t.Run(
		"uint8",
		func(t *testing.T) {
			value = uint8(1)
			expect = []byte{1}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}

			value = lifxlan.AckResFlag(2)
			expect = []byte{2}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}

			value = lifxlan.ServiceType(3)
			expect = []byte{3}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}
		},
	)

	t.Run(
		"uint16",
		func(t *testing.T) {
			value = uint16(0x0102)
			expect = []byte{2, 1}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}

			value = lifxlan.TaggedHeader(0x0201)
			expect = []byte{1, 2}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}

			value = lifxlan.MessageType(0x0100)
			expect = []byte{0, 1}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}
		},
	)

	t.Run(
		"uint32",
		func(t *testing.T) {
			value = uint32(0x01020304)
			expect = []byte{4, 3, 2, 1}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}
		},
	)

	t.Run(
		"uint64",
		func(t *testing.T) {
			value = uint64(0x0102030405060708)
			expect = []byte{8, 7, 6, 5, 4, 3, 2, 1}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}

			value = lifxlan.Target(0x0706050403020100)
			expect = []byte{0, 1, 2, 3, 4, 5, 6, 7}
			actual = lifxlan.UintToBytes(value)
			if !bytes.Equal(actual, expect) {
				t.Errorf(
					"UintToBytes(%v) expected % x, got % x",
					value,
					expect,
					actual,
				)
			}
		},
	)
}
