package auth

type Service interface {
	FindAuth(authD *AuthDetails) bool
	FetchAuth(*AuthDetails) (*Auth, error)
	DeleteAuth(*AuthDetails) error
	CreateAuth(uint64) (*Auth, error)
	Login(AuthDetails) (string, error)
	FindAuthUser(userID int) bool
	DeleteAuthUser(userID int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAuth(authD *AuthDetails) bool {
	err := s.repository.FindAuth(authD)
	if err {
		return true
	} else {
		return false
	}
}

func (s *service) FetchAuth(authD *AuthDetails) (*Auth, error) {
	au, err := s.repository.FetchAuth(authD)
	if err != nil {
		return nil, err
	}

	return au, nil

}

func (s *service) DeleteAuth(authD *AuthDetails) error {
	err := s.repository.DeleteAuth(authD)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CreateAuth(userID uint64) (*Auth, error) {
	au, err := s.repository.CreateAuth(userID)
	if err != nil {
		return nil, err
	}
	return au, nil
}

func (s *service) Login(authD AuthDetails) (string, error) {
	token, err := CreateToken(authD)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *service) FindAuthUser(userID int) bool {
	err := s.repository.FindAuthUser(userID)
	if err {
		return true
	} else {
		return false
	}
}

func (s *service) DeleteAuthUser(userID int) error {
	err := s.repository.DeleteAuthUser(userID)
	if err != nil {
		return err
	}
	return nil
}
