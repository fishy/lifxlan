package lifxlan_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"go.yhsif.com/lifxlan"
)

func features2json(t *testing.T, features lifxlan.Features) string {
	t.Helper()

	var sb strings.Builder
	if err := json.NewEncoder(&sb).Encode(features); err != nil {
		t.Errorf("Failed to encode %#v to json: %v", features, err)
	}
	return strings.TrimSpace(sb.String())
}

func compareFeatures(t *testing.T, expected, actual lifxlan.Features) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, actual %v", features2json(t, expected), features2json(t, actual))
	}
}

func TestFeatureAt(t *testing.T) {
	const (
		major  = 1
		minor1 = 1
		minor2 = 2
	)
	baseFeature := lifxlan.Features{
		HEV:   lifxlan.OptionalBoolPtr(false),
		Color: lifxlan.OptionalBoolPtr(false),
	}
	upgradeFeature1 := lifxlan.Features{
		HEV: lifxlan.OptionalBoolPtr(true),
	}
	upgradeFeature2 := lifxlan.Features{
		Color: lifxlan.OptionalBoolPtr(true),
	}

	expectedUpgrade1 := lifxlan.Features{
		HEV:   lifxlan.OptionalBoolPtr(true),
		Color: lifxlan.OptionalBoolPtr(false),
	}
	expectedUpgrade2 := lifxlan.Features{
		HEV:   lifxlan.OptionalBoolPtr(true),
		Color: lifxlan.OptionalBoolPtr(true),
	}

	product := lifxlan.Product{
		Features: baseFeature,
		Upgrades: []lifxlan.FirmwareUpgrade{
			{
				Major:    major,
				Minor:    minor1,
				Features: upgradeFeature1,
			},
			{
				Major:    major,
				Minor:    minor2,
				Features: upgradeFeature2,
			},
		},
	}

	compareFeatures(t, expectedUpgrade1, product.FeaturesAt(lifxlan.FirmwareUpgrade{
		Major: major,
		Minor: minor1,
	}))
	compareFeatures(t, expectedUpgrade2, product.FeaturesAt(lifxlan.FirmwareUpgrade{
		Major: major,
		Minor: minor2,
	}))
}

func TestFirmwareUpgradeLess(t *testing.T) {
	versions := []lifxlan.FirmwareUpgrade{
		{
			Major: 1,
			Minor: 1,
		},
		{
			Major: 1,
			Minor: 2,
		},
		{
			Major: 1,
			Minor: 10,
		},
		{
			Major: 2,
			Minor: 1,
		},
		{
			Major: 2,
			Minor: 2,
		},
		{
			Major: 2,
			Minor: 10,
		},
		{
			Major: 10,
			Minor: 1,
		},
		{
			Major: 10,
			Minor: 2,
		},
		{
			Major: 10,
			Minor: 10,
		},
	}

	t.Run("equal", func(t *testing.T) {
		for _, v := range versions {
			if v.Less(v) {
				t.Errorf("Expected %v.Less(%v) to be false, got true", v, v)
			}
		}
	})

	t.Run("less", func(t *testing.T) {
		for i, v1 := range versions {
			for j := i + 1; j < len(versions); j++ {
				v2 := versions[j]
				if !v1.Less(v2) {
					t.Errorf("Expected %v.Less(%v) to be true, got false", v1, v2)
				}
				if v2.Less(v1) {
					t.Errorf("Expected %v.Less(%v) to be false, got true", v2, v1)
				}
			}
		}
	})
}
