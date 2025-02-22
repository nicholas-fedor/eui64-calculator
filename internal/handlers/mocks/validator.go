package mocks

import "errors"

var (
	ErrInvalidMAC    = errors.New("invalid MAC")
	ErrInvalidPrefix = errors.New("invalid prefix")
)

type Validator struct {
	MacErr    error
	PrefixErr error
}

func (m *Validator) ValidateMAC(_ string) error {
	return m.MacErr
}

func (m *Validator) ValidateIPv6Prefix(_ string) error {
	return m.PrefixErr
}
