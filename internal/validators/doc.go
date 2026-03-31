// Package validators provides input validation functions for the EUI-64 calculator.
//
// It validates MAC addresses and IPv6 prefixes before computing EUI-64 interface
// identifiers. Each validation function returns a sentinel error on failure that
// describes the specific validation rule violated.
//
// The package exports two primary validation functions:
//   - ValidateMAC: validates a 48-bit MAC address string
//   - ValidateIPv6Prefix: validates an IPv6 network prefix string (first 64 bits)
//
// All exported error variables follow Go conventions for sentinel errors and can
// be checked using errors.Is().
//
// Tests in *_test.go files provide table-driven coverage for both validators,
// including edge cases such as empty input, overflow, and invalid characters.
package validators
