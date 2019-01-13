package light

import (
	"context"
	"math"
	"net"
	"time"

	"github.com/fishy/lifxlan"
)

// Waveform defines the type of the waveform.
//
// https://lan.developer.lifx.com/v2.0/docs/waveforms
type Waveform uint8

// Waveform values.
//
// https://lan.developer.lifx.com/v2.0/docs/waveforms
const (
	// https://lan.developer.lifx.com/v2.0/docs/waveforms#section-saw
	WaveformSaw Waveform = 0

	// https://lan.developer.lifx.com/v2.0/docs/waveforms#section-sine
	WaveformSine Waveform = 1

	// https://lan.developer.lifx.com/v2.0/docs/waveforms#section-half-sine
	WaveformHalfSine Waveform = 2

	// https://lan.developer.lifx.com/v2.0/docs/waveforms#section-triangle
	WaveformTriangle Waveform = 3

	// https://lan.developer.lifx.com/v2.0/docs/waveforms#section-pulse
	WaveformPulse Waveform = 4
)

// BoolUint8 is the uint8 value used to represent bool.
type BoolUint8 uint8

// Bool2Uint8 converts bool value into BoolUint8 value.
func Bool2Uint8(b bool) BoolUint8 {
	if b {
		return 1
	}
	return 0
}

// ConvertSkewRatio scales [0, 1] into [-32768, 32767].
func ConvertSkewRatio(v float64) int16 {
	return int16(int64(math.Round(v*float64(math.MaxUint16))) - 32768)
}

// RawSetWaveformOptionalPayload defines the struct to be used for encoding and
// decoding.
//
// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-setwaveformoptional-119
type RawSetWaveformOptionalPayload struct {
	_             uint8 // reserved
	Transient     BoolUint8
	Color         lifxlan.Color
	Period        lifxlan.TransitionTime
	Cycles        float32
	SkewRatio     int16
	Waveform      Waveform
	SetHue        BoolUint8
	SetSaturation BoolUint8
	SetBrightness BoolUint8
	SetKelvin     BoolUint8
}

// SetWaveformArgs is the args to be translated into
// RawSetWaveformOptionalPayload.
type SetWaveformArgs struct {
	// True means that after the waveform it should go back to its original color.
	Transient bool

	// The target color.
	Color *lifxlan.Color

	// Duratino of a cycle.
	Period time.Duration

	// Number of cycles.
	Cycles float32

	// Type of waveform.
	Waveform Waveform

	// SkewRatio should be in range [0, 1] and it is only used with WaveformPulse.
	//
	// https://lan.developer.lifx.com/v2.0/docs/waveforms#section-pulse
	SkewRatio float64

	// The color args with Keep* set to true will not be changed.
	//
	// For example the current color has H 255, S 255, B 255 and K 255,
	// and the target color has H 0, S 0, B 0, K 0.
	// When KeepHue is false and all else are true,
	// Only H will be changed during the S, B, and K won't be changed during the
	// waveform.
	//
	// Please note this is the reverse of the set_* definitions of
	// RawSetWaveformOptionalPayload.
	// The reason is to make sure that when they are all zero values,
	// it behaves the same as SetWaveform message as defined in:
	// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-setwaveform-103
	KeepHue        bool
	KeepSaturation bool
	KeepBrightness bool
	KeepKelvin     bool
}

func (ld *device) SetWaveform(
	ctx context.Context,
	conn net.Conn,
	args *SetWaveformArgs,
	ack bool,
) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := ld.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	var flags lifxlan.AckResFlag
	if ack {
		flags |= lifxlan.FlagAckRequired
	}

	// Send
	seq, err := ld.Send(
		ctx,
		conn,
		flags,
		SetWaveformOptional,
		&RawSetWaveformOptionalPayload{
			Transient:     Bool2Uint8(args.Transient),
			Color:         ld.SanitizeColor(*args.Color),
			Period:        lifxlan.ConvertDuration(args.Period),
			Cycles:        args.Cycles,
			SkewRatio:     ConvertSkewRatio(1 - args.SkewRatio),
			Waveform:      args.Waveform,
			SetHue:        Bool2Uint8(!args.KeepHue),
			SetSaturation: Bool2Uint8(!args.KeepSaturation),
			SetBrightness: Bool2Uint8(!args.KeepBrightness),
			SetKelvin:     Bool2Uint8(!args.KeepKelvin),
		},
	)
	if err != nil {
		return err
	}

	if ack {
		return lifxlan.WaitForAcks(ctx, conn, ld.Source(), seq)
	}
	return nil
}
