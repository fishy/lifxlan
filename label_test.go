package lifxlan_test

import (
	"context"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/mock"
)

func TestLabel(t *testing.T) {
	t.Run(
		"Short",
		func(t *testing.T) {
			label := "1234"
			var rl lifxlan.Label
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("Label.Set returned err: %v", err)
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
			var rl lifxlan.Label
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("Label.Set returned err: %v", err)
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
			var rl lifxlan.Label
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("Label.Set returned err: %v", err)
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
			var rl lifxlan.Label
			err := rl.Set(label)
			if err != nil {
				t.Fatalf("Label.Set returned err: %v", err)
			}
			got := rl.String()
			if got != expected {
				t.Errorf("Expected %q, got %q", label, got)
			}
		},
	)
}

func TestEmptyLabel(t *testing.T) {
	var label lifxlan.Label
	s := label.String()
	if s != lifxlan.EmptyLabel {
		t.Errorf("Expected %q, got %q", lifxlan.EmptyLabel, s)
	}
}

func TestGetLabel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	var expected lifxlan.Label
	expected.Set("foo")

	service, device := mock.StartService(t)
	service.RawStateLabelPayload = &lifxlan.RawStateLabelPayload{
		Label: expected,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := device.GetLabel(ctx, nil); err != nil {
		t.Fatal(err)
	}
	if device.Label().String() != expected.String() {
		t.Errorf("Label expected %v, got %v", expected, device.Label())
	}
}
