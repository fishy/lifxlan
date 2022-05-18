package lifxlan

import (
	"fmt"
	"sort"
)

// OptionalBool defines a helper type for optional boolean fields
type OptionalBool bool

// OptionalBoolPtr is a helper function for writing *OptionalBool literal.
func OptionalBoolPtr(v OptionalBool) *OptionalBool {
	return &v
}

// Get returns false when ob is unset, otherwise it returns the set value.
func (ob *OptionalBool) Get() bool {
	if ob == nil {
		return false
	}
	return bool(*ob)
}

// Fallback returns, in this order:
//
// - ob, if it's set,
//
// - a copy of fallback if fallback is set,
//
// - nil if neither ob nor fallback is set.
func (ob *OptionalBool) Fallback(fallback *OptionalBool) *OptionalBool {
	if ob != nil {
		return ob
	}
	if fallback != nil {
		v := *fallback
		return &v
	}
	return nil
}

// TemperatureRange defines the json format of temperature range of a product.
//
// It would be either a slice of length 0 (meaning this is not a light device),
// or length 2 (min, max).
type TemperatureRange []uint16

// The magic length of TemperatureRange to be considered as valid (min, max).
const (
	ValidTemperatureRangeLength = 2
)

// Fallback returns, in this order:
//
// - tr if it's valid,
//
// - a copy of fallback if it's valid,
//
// - nil if neither tr nor fallback is valid.
func (tr TemperatureRange) Fallback(fallback TemperatureRange) TemperatureRange {
	if tr.Valid() {
		return tr
	}
	if fallback.Valid() {
		v := make(TemperatureRange, ValidTemperatureRangeLength)
		copy(v, fallback)
		return v
	}
	return nil
}

// Valid returns true if tr has a length of exactly 2.
func (tr TemperatureRange) Valid() bool {
	return len(tr) == ValidTemperatureRangeLength
}

// Min returns the min temperature if tr is valid, 0 otherwise.
func (tr TemperatureRange) Min() uint16 {
	if tr.Valid() {
		return tr[0]
	}
	return 0
}

// Max returns the max temperature if tr is valid, 0 otherwise.
func (tr TemperatureRange) Max() uint16 {
	if tr.Valid() {
		return tr[1]
	}
	return 0
}

// Features defines the json format of features of a product.
type Features struct {
	HEV               *OptionalBool `json:"hev,omitempty"`
	Color             *OptionalBool `json:"color,omitempty"`
	Chain             *OptionalBool `json:"chain,omitempty"`
	Matrix            *OptionalBool `json:"matrix,omitempty"`
	Relays            *OptionalBool `json:"relays,omitempty"`
	Buttons           *OptionalBool `json:"buttons,omitempty"`
	Infrared          *OptionalBool `json:"infrared,omitempty"`
	Multizone         *OptionalBool `json:"multizone,omitempty"`
	ExtendedMultizone *OptionalBool `json:"extended_multizone,omitempty"`

	TemperatureRange TemperatureRange `json:"temperature_range,omitempty"`
}

// MergeFeatures merges the features defined in features,
// Each feature falls back to the next one in features if it's unset.
func MergeFeatures(features ...Features) Features {
	var result Features
	for _, f := range features {
		result.HEV = result.HEV.Fallback(f.HEV)
		result.Color = result.Color.Fallback(f.Color)
		result.Chain = result.Chain.Fallback(f.Chain)
		result.Matrix = result.Matrix.Fallback(f.Matrix)
		result.Relays = result.Relays.Fallback(f.Relays)
		result.Buttons = result.Buttons.Fallback(f.Buttons)
		result.Infrared = result.Infrared.Fallback(f.Infrared)
		result.Multizone = result.Multizone.Fallback(f.Multizone)
		result.ExtendedMultizone = result.ExtendedMultizone.Fallback(f.ExtendedMultizone)

		result.TemperatureRange = result.TemperatureRange.Fallback(f.TemperatureRange)
	}
	return result
}

// Vendor defines a vendor.
type Vendor struct {
	ID       uint32
	Name     string
	Defaults Features
}

// FirmwareUpgrade defines a firmware version with optional upgrade features.
type FirmwareUpgrade struct {
	Major    uint16   `json:"major"`
	Minor    uint16   `json:"minor"`
	Features Features `json:"features"`
}

// Less returns true if fu's firmware version is smaller than other's firmware
// version.
func (fu FirmwareUpgrade) Less(other FirmwareUpgrade) bool {
	if fu.Major < other.Major {
		return true
	}
	if fu.Major > other.Major {
		return false
	}
	return fu.Minor < other.Minor
}

func (fu FirmwareUpgrade) String() string {
	return fmt.Sprintf("(%d, %d)", fu.Major, fu.Minor)
}

// EmptyFirmware is the constant to be compared against
// Device.Firmware().String().
const EmptyFirmware = "(0, 0)"

// Upgrades defines sortable interface of FirmwareUpgrade.
type Upgrades []FirmwareUpgrade

var _ sort.Interface = Upgrades(nil)

func (u Upgrades) Len() int           { return len(u) }
func (u Upgrades) Less(i, j int) bool { return u[i].Less(u[j]) }
func (u Upgrades) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

// Product defines a product.
type Product struct {
	VendorID   uint32 `json:"-"`
	VendorName string `json:"-"`

	ProductID   uint32   `json:"pid"`
	ProductName string   `json:"name"`
	Features    Features `json:"features"`
	Upgrades    Upgrades `json:"upgrades"`
}

// FeaturesAt gets the features at the given firmware version,
// with appropriate upgrades applied.
//
// Calling with zero firmware will return the same features as the Features
// field.
func (p Product) FeaturesAt(firmware FirmwareUpgrade) Features {
	features := make([]Features, 0, len(p.Upgrades)+1)
	for _, u := range p.Upgrades {
		if firmware.Less(u) {
			continue
		}
		features = append(features, u.Features)
	}
	features = append(features, p.Features)
	return MergeFeatures(features...)
}
