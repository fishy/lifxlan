package main

import (
	"context"
	"log"
	"os"

	"github.com/fishy/lifxlan"
)

func checkContextError(err error) bool {
	return err != nil && err != context.Canceled && err != context.DeadlineExceeded
}

func main() {
	log.Print("Reading image from stdin...")
	img, err := readImage(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	td := findDevice(lifxlan.Target(target))
	if td == nil {
		log.Fatal("No matching tile device found.")
	}
	log.Printf("Found %v", td)
	draw(td, img)
}
