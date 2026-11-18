//go:build !js || !wasm

package main

import "fmt"

func main() {
	fmt.Println("WebDataEditer4Go CLI / Native Stub. For WebAssembly, build with GOOS=js GOARCH=wasm.")
}
