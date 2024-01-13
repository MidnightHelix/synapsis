package service

type UserService interface {
	Register(username, password string) error
	Login(username, password string) (bool, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(username, password string) error {
	// You can add validation logic here before registering the user
	return s.repo.Register(username, password)
}

func (s *userService) Login(username, password string) (bool, error) {
	return s.repo.Login(username, password)
}
