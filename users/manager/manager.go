package manager

type UsersStorage interface {
}

type Manager struct {
}

type Email string

func New(s UsersStorage) *Manager {
	return nil
}

func (m *Manager) CreateUser(email Email, fullname string, password string) (string, error) {
	return "", nil
}
