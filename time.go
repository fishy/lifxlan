package lifxlan

import (
	"time"
)

// msPerSecond is the number of milliseconds per second (1000).
//
// It's of the type time.Duration so in most cases you need to cast it.
const msPerSecond = time.Second / time.Millisecond

// Timestamp is the type used in messages to represent a timestamp.
//
// It's defined as milliseconds since UNIX EPOCH.
type Timestamp uint64

// ConvertTime converts a time.Time into Timestamp.
func ConvertTime(t time.Time) Timestamp {
	return Timestamp(
		t.Unix()*int64(msPerSecond) + int64(t.Nanosecond())/int64(time.Millisecond),
	)
}

// Time converts a Timestamp into time.Time.
func (ts Timestamp) Time() time.Time {
	sec := uint64(ts) / uint64(msPerSecond)
	nano := uint64(ts) % uint64(msPerSecond) * uint64(time.Millisecond)
	return time.Unix(int64(sec), int64(nano))
}

// TransitionTime is the type used in messages to represent transition time.
//
// Its unit is milliseconds.
type TransitionTime uint32

// ConvertDuration converts a time.Duration into TransitionTime.
//
// The max uint32 value can represent a transition time of more than 1,193
// hours[1] (or, in other words, more than a month).
// So although an overflow is technically possible,
// we don't really do any special handlings here
// (it's not a security risk and won't crash anything[2]).
// If you feed in a duration that overflows TransitionTime,
// you should feel bad (or great, it's totally up to you) about it.
// Do you really want your light(s) to take more than a month to turn on/off?
//
// [1] https://play.golang.com/p/LqfMpvhIctx
//
// [2] https://play.golang.com/p/edwqG4nNqkt
func ConvertDuration(d time.Duration) TransitionTime {
	return TransitionTime(d / time.Millisecond)
}

// Duration converts a TransitionTime into time.Duration.
func (tt TransitionTime) Duration() time.Duration {
	return time.Duration(time.Duration(tt) * time.Millisecond)
}
