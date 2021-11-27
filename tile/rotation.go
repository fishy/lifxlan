package tile

import (
	"fmt"
	"math"
)

// Rotation defines the rotation of a single tile.
//
// NOTE: Currently only RotationRightSideUp is fully supported.
type Rotation int

// Possible Rotation values.
const (
	RotationRightSideUp Rotation = iota
	RotationRotateRight
	RotationRotateLeft
	RotationFaceDown
	RotationFaceUp
	RotationUpsideDown
)

func (r Rotation) String() string {
	switch r {
	default:
		return fmt.Sprintf("<Invalid value %d>", int(r))
	case RotationRightSideUp:
		return "RightSideUp"
	case RotationRotateRight:
		return "RotateRight"
	case RotationRotateLeft:
		return "RotateLeft"
	case RotationFaceDown:
		return "FaceDown"
	case RotationFaceUp:
		return "FaceUp"
	case RotationUpsideDown:
		return "UpsideDown"
	}
}

// ParseRotation parses accelerator measurements into Rotation
func ParseRotation(x, y, z int16) Rotation {
	abs := func(x int16) int16 {
		return int16(math.Abs(float64(x)))
	}

	// Copied from:
	// https://lan.developer.lifx.com/docs/tile-control#tile-orientation
	absX := abs(x)
	absY := abs(y)
	absZ := abs(z)

	if x == -1 && y == -1 && z == -1 {
		// Invalid data, assume right-side up.
		return RotationRightSideUp
	}
	if absX > absY && absX > absZ {
		if x > 0 {
			return RotationRotateRight
		}
		return RotationRotateLeft
	}

	if absZ > absX && absZ > absY {
		if z > 0 {
			return RotationFaceDown
		}
		return RotationFaceUp
	}

	if y > 0 {
		return RotationUpsideDown
	}
	return RotationRightSideUp
}
