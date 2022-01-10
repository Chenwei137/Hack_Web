package main

import "syscall/js"

func main() {
	hi_wasm := js.Global().Get("hi_wasm")
	hi_wasm.Invoke("WebAssembly")
}
