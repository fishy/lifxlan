package lifxlan

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Target defines a target by its MAC address.
type Target uint64

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
