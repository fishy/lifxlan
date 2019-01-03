package lifxlan

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"strings"
)

// Target defines a target by its MAC address.
type Target uint64

var _ flag.Value = (*Target)(nil)

// AllDevices is the special Target value means all devices.
const AllDevices Target = 0

func (t Target) String() string {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(t))
	strs := make([]string, 6)
	for i, b := range buf[:6] {
		strs[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(strs, ":")
}

// Set implements flag.Value interface.
func (t *Target) Set(s string) (err error) {
	*t, err = ParseTarget(s)
	return
}

// Matches returns true if either target is AllDevices,
// or both targets have the same value.
func (t Target) Matches(other Target) bool {
	if t == other {
		return true
	}
	if t == AllDevices {
		return true
	}
	if other == AllDevices {
		return true
	}
	return false
}

// ParseTarget parses s into a Target.
//
// s should be in the format of a MAC address, e.g. "01:23:45:67:89:ab",
// or the special value for AllDevices: "00:00:00:00:00:00" or "".
func ParseTarget(s string) (t Target, err error) {
	// Special case.
	if s == "" {
		return AllDevices, nil
	}

	var mac net.HardwareAddr
	mac, err = net.ParseMAC(s)
	if err != nil {
		return
	}
	buf := make([]byte, 8)
	copy(buf, mac)
	t = Target(binary.LittleEndian.Uint64(buf))
	return
}
