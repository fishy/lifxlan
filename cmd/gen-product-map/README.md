# gen-product-map

This is the helper tool to generate
[`ProductMap`](https://pkg.go.dev/go.yhsif.com/lifxlan#pkg-variables)
defined in
[`product_map.go`](https://github.com/fishy/lifxlan/blob/master/product_map.go).

## Installation

```
go get -u go.yhsif.com/lifxlan/cmd/gen-product-map
```

## Usage

```sh
gen-product-map >> product_map.go
# Then manally update the file to remove previous value.
```
