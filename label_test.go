package lifxlan_test

import (
	"testing"

	"github.com/fishy/lifxlan"
)

func TestRawLabel(t *testing.T) {
	t.Run(
		"Short",
		func(t *testing.T) {
			label := "1234"
			var rl lifxlan.RawLabel
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("RawLabel.Set returned err: %v", err)
			}
			got := rl.String()
			if got != label {
				t.Errorf("Expected %q, got %q", label, got)
			}
		},
	)

	t.Run(
		"UTF-8",
		func(t *testing.T) {
			label := "中文"
			var rl lifxlan.RawLabel
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("RawLabel.Set returned err: %v", err)
			}
			got := rl.String()
			if got != label {
				t.Errorf("Expected %q, got %q", label, got)
			}
		},
	)

	t.Run(
		"Long",
		func(t *testing.T) {
			label := "0123456789012345678901234567890123456789"
			// First 32 bytes in utf8
			expected := "01234567890123456789012345678901"
			var rl lifxlan.RawLabel
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("RawLabel.Set returned err: %v", err)
			}
			got := rl.String()
			if got != expected {
				t.Errorf("Expected %q, got %q", label, got)
			}
		},
	)

	t.Run(
		"LongUnicode",
		func(t *testing.T) {
			label := "中文6789012345678901234567890123456789"
			// First 32 bytes in utf8
			expected := "中文67890123456789012345678901"
			var rl lifxlan.RawLabel
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("RawLabel.Set returned err: %v", err)
			}
			got := rl.String()
			if got != expected {
				t.Errorf("Expected %q, got %q", label, got)
			}
		},
	)
}
