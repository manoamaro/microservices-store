package mock

type AuthService struct {
}

func (m *AuthService) Validate(token string, audiences ...string) (bool, error) {
	return true, nil
}
