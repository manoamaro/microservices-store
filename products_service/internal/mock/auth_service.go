package mock

type MockAuthService struct {
}

func (m *MockAuthService) Validate(token string, audiences ...string) (error, bool) {
	return nil, true
}
