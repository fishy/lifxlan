[![GoDoc](https://godoc.org/github.com/fishy/lifxlan?status.svg)](https://pkg.go.dev/github.com/fishy/lifxlan)
[![Go Report Card](https://goreportcard.com/badge/github.com/fishy/lifxlan)](https://goreportcard.com/report/github.com/fishy/lifxlan)

# LIFX LAN

This is a library that provides API implemented in Go for
[LIFX LAN Protocol](https://lan.developer.lifx.com/v2.0/docs/).

## Overview

The root package focuses on the base stuff, device discovery,
and capabilities shared by all types of devices.
Subpackages provide more concreted capabilities by different types of LIFX
devices like light control and tile control.

Currently this library is not complete and implement all possible LIFX LAN
Protocols is not the current goal of this library.
The design choice for this library is that it exposes as much as possible,
so another third party package can implement missing device APIs by wrapping
`Device`s returned by this package.
Please refer to the subpackage code for an example of extending device
capabilities.
The reason its split into subpackages is to make sure that it's extendible.

The main focus right now is on
[tile API](https://lan.developer.lifx.com/v2.0/docs/tile-control) support.
The reason is that at the time of writing,
although there are several Go projects implemented LIFX LAN Protocol available,
none of them support tile APIs.
Please refer to
[`tile` subpackage on GoDoc](https://pkg.go.dev/github.com/fishy/lifxlan/tile)
for more details.

All API with (potential) I/O calls takes a [context](https://pkg.go.dev/context)
arg and checks for (and in most cases, relies on) context cancellations.

The API is unstable right now, but I try very hard not to break them.

## Examples

Besides
[examples on GoDoc](https://pkg.go.dev/github.com/fishy/lifxlan#pkg-examples),
there are also some example command line apps in
[`lifxlan-examples`](https://github.com/fishy/lifxlan-examples) repository.

## License

[BSD License](https://github.com/fishy/lifxlan/blob/master/LICENSE).
