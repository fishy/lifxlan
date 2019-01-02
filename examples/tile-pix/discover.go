package main

import (
	"context"
	"log"
	"sync"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/tile"
)

func findDevice(target lifxlan.Target) (td tile.Device) {
	var ctx context.Context
	var cancel context.CancelFunc
	if *discoverTimeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), *discoverTimeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
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

	for device := range deviceChan {
		if !device.Target().Matches(target) {
			continue
		}

		wg.Add(1)
		go func(device lifxlan.Device) {
			defer wg.Done()
			log.Printf("Found %v, checking tile capablities...", device)
			t, err := tile.Wrap(ctx, device, false)
			if checkContextError(err) {
				log.Printf("Check tile capablities for %v failed: %v", device, err)
			} else {
				if t == nil {
					return
				}
				td = t
				cancel()
			}
		}(device)
	}

	wg.Wait()
	return
}
