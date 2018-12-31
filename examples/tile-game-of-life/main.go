package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/fishy/lifxlan"
)

// Flags
var (
	timeout = flag.Duration(
		"discoverTimeout",
		time.Second*2,
		"Timeout for discover API calls",
	)

	targetStr = flag.String(
		"target",
		"00:00:00:00:00:00",
		"The MAC address of the target tile device. Default value means any (first) tile device",
	)

	broadcastHost = flag.String(
		"broadcastHost",
		"",
		`Broadcast IP (e.g. "192.168.1.255"). Empty value means "255.255.255.255", which should work in most networks`,
	)
)

func checkContextError(err error) bool {
	return err != nil && err != context.Canceled && err != context.DeadlineExceeded
}

func main() {
	flag.Parse()

	target, err := lifxlan.ParseTarget(*targetStr)
	if err != nil {
		log.Fatalf("Unable to parse target %q: %v", *targetStr, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	deviceChan := make(chan lifxlan.Device)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := lifxlan.Discover(ctx, deviceChan, ""); err != nil {
			if checkContextError(err) {
				log.Print("Discover failed:", err)
			}
		}
	}()

	var tile *lifxlan.TileDevice

	for device := range deviceChan {
		if !device.Target().Matches(target) {
			continue
		}

		wg.Add(1)
		go func(device lifxlan.Device) {
			defer wg.Done()
			t, err := device.GetTileDevice(ctx)
			if checkContextError(err) {
				log.Printf("Check tile for %v failed: %v\n", device, err)
			} else {
				if t == nil {
					return
				}
				tile = t
				cancel()
			}
		}(device)
	}

	wg.Wait()
	log.Print(tile)
}
