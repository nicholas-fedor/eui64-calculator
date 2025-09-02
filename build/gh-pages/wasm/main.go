//go:build js && wasm
// +build js,wasm

// Package main provides a WebAssembly module for client-side EUI-64 calculations.
// It exposes functions to validate MAC addresses, IPv6 prefixes, and compute EUI-64
// identifiers, integrating with the browser's JavaScript environment.
package main

import (
	"syscall/js"

	"github.com/nicholas-fedor/eui64-calculator/internal/eui64"
	"github.com/nicholas-fedor/eui64-calculator/internal/validators"
)

// main initializes the WebAssembly module, registering JavaScript functions and
// keeping the module alive in the browser event loop.
func main() {
	js.Global().Set("validateMAC", js.FuncOf(validateMACFunc))
	js.Global().Set("validateIPv6Prefix", js.FuncOf(validateIPv6PrefixFunc))
	js.Global().Set("calculateEUI64", js.FuncOf(calculateEUI64Func))
	<-make(chan bool) // Block indefinitely to keep WASM module active.
}

// validateMACFunc validates a MAC address string provided via JavaScript.
// It expects a single string argument and returns an empty string on success or
// an error message on failure.
func validateMACFunc(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "Invalid number of arguments"
	}
	mac := args[0].String()
	if err := validators.ValidateMAC(mac); err != nil {
		return err.Error()
	}
	return ""
}

// validateIPv6PrefixFunc validates an IPv6 prefix string provided via JavaScript.
// It expects a single string argument and returns an empty string on success or
// an error message on failure.
func validateIPv6PrefixFunc(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "Invalid number of arguments"
	}
	prefix := args[0].String()
	if err := validators.ValidateIPv6Prefix(prefix); err != nil {
		return err.Error()
	}
	return ""
}

// calculateEUI64Func computes the EUI-64 interface ID and full IPv6 address from
// a MAC address and IPv6 prefix provided via JavaScript. It expects two string
// arguments (MAC and prefix) and returns a JavaScript object with "interfaceID"
// and "fullIP" fields on success, or an error message on failure.
func calculateEUI64Func(this js.Value, args []js.Value) any {
	if len(args) != 2 {
		return "Invalid number of arguments"
	}
	mac := args[0].String()
	prefix := args[1].String()
	interfaceID, fullIP, err := eui64.CalculateEUI64(mac, prefix)
	if err != nil {
		return err.Error()
	}
	return js.ValueOf(map[string]interface{}{
		"interfaceID": interfaceID,
		"fullIP":      fullIP,
	})
}
