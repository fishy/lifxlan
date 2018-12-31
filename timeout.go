package lifxlan

import (
	"time"
)

// UDPReadTimeout is the read timeout we use to read all the UDP messages.
//
// In some functions (e.g. Discover),
// the function will simply use the timeout to check context cancellation and
// continue reading,
// instead of return upon timeout.
//
// It's intentionally defined as variable instead of constant,
// so the user could adjust it if needed.
var UDPReadTimeout = time.Millisecond * 100

// GetReadDeadline returns a value can be used in net.Conn.SetReadDeadline from
// UDPReadTimeout value.
func GetReadDeadline() time.Time {
	return time.Now().Add(UDPReadTimeout)
}

type timeouter interface {
	Timeout() bool
}

// CheckTimeoutError returns true if err is caused by timeout in net package.
func CheckTimeoutError(err error) bool {
	if t, ok := err.(timeouter); ok {
		return t.Timeout()
	}
	return false
}
