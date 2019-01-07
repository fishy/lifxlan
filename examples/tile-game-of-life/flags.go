package main

import (
	"errors"
	"flag"
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/fishy/lifxlan"
)

var (
	noack = flag.Bool(
		"noack",
		false,
		"Do not require ack for drawing API calls.",
	)

	discoverTimeout = flag.Duration(
		"discoverTimeout",
		time.Second*2,
		"Timeout for discover API calls.",
	)

	drawTimeout = flag.Duration(
		"drawTimeout",
		time.Millisecond*200,
		"Timeout for drawing API calls.",
	)

	interval = flag.Duration(
		"interval",
		time.Millisecond*1500,
		"Interval between 2 frames.",
	)

	broadcastHost = flag.String(
		"broadcastHost",
		"",
		`Broadcast IP (e.g. "192.168.1.255"). Empty value means "255.255.255.255", which should work in most networks.`,
	)

	generations = flag.Int(
		"generations",
		20,
		"Number of generations before stop. 0 means never stop (but still stops when the board is empty).",
	)

	loop = flag.Bool(
		"loop",
		false,
		"After reached the number of generations, reset the board and loop over instead of reverting to the original colors.",
	)

	kelvin = flag.Int(
		"kelvin",
		4000,
		"The Kelvin value of the color, in range of [2500, 9000].",
	)

	turnon = flag.Bool(
		"turnon",
		false,
		"Turn on the device if it's not already on, and turn it off afterwards (when no loop args specified).",
	)

	origColor = flagColor{0xff, 0xff, 0xff}

	target lifxlan.Target
)

type flagColor [3]uint8

func (c flagColor) String() string {
	return fmt.Sprintf("%x", []uint8(c[:]))
}

func (c *flagColor) Set(s string) error {
	v, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return err
	}
	if v > 0xffffff {
		return errors.New("value out of range")
	}
	(*c)[0] = uint8((v & 0xff0000) >> 16)
	(*c)[1] = uint8((v & 0xff00) >> 8)
	(*c)[2] = uint8(v & 0xff)
	return nil
}

func (c flagColor) RGBA() (uint32, uint32, uint32, uint32) {
	return color.RGBA{
		R: c[0],
		G: c[1],
		B: c[2],
	}.RGBA()
}

func init() {
	flag.Var(
		&target,
		"target",
		"The MAC address of the target tile device. Empty value means any (first) tile device.",
	)
	flag.Var(
		&origColor,
		"color",
		`The hex color to use, in format of "rrggbb".`,
	)
	flag.Parse()
}
