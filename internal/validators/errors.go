package validators

import "errors"

// Exported errors for validator operations.
var (
	ErrMACRequired         = errors.New("MAC address is required")
	ErrInvalidMACLength    = errors.New("invalid MAC address length")
	ErrParseMACFailed      = errors.New("failed to parse MAC address")
	ErrPrefixRequired      = errors.New("a non-blank IPv6 prefix is expected")
	ErrPrefixTooLong       = errors.New("IPv6 prefix exceeds maximum hextets")
	ErrEmptyHextet         = errors.New("empty hextet in IPv6 prefix")
	ErrInvalidHextetChar   = errors.New("invalid character in hextet")
	ErrInvalidHextetLength = errors.New("invalid hextet length in IPv6 prefix")
)
