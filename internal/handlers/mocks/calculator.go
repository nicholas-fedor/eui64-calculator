package mocks

type Calculator struct {
	Err         error
	InterfaceID string
	FullIP      string
}

func (m *Calculator) CalculateEUI64(_, _ string) (string, string, error) {
	return m.InterfaceID, m.FullIP, m.Err
}
