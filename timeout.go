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

type timeouter interface {
	Timeout() bool
}

func getReadDeadline() time.Time {
	return time.Now().Add(UDPReadTimeout)
}
