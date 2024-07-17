package app

import (
	validator "github.com/YuliaKorolenko/valid"
	"homework9/internal/entities/user"
)

type URepository interface {
	CreateUser(user user.User) (user.User, error)
	UpdateUser(user user.User) (user.User, error)
	GetCurrentUser(id int64) (user.User, error)
	DeleteUser(user user.User) error
}

func (a App) CreateUser(nickname string, email string) (user.User, error) {
	u := user.NewUser(-1, nickname, email)
	err := validator.Validate(u)
	if err != nil {
		return u, ValidationApError{ErrorValidation}
	}
	return a.userRepository.CreateUser(u)
}

func (a App) UpdateUser(id int64, nickname string, email string) (user.User, error) {
	u, err := a.userRepository.GetCurrentUser(id)
	if err != nil {
		return u, NoSuchUserError{ErrorNoSuchUser}
	}
	u = user.NewUser(u.ID, nickname, email)

	err = validator.Validate(u)
	if err != nil {
		return u, ValidationApError{ErrorValidation}
	}
	return a.userRepository.UpdateUser(u)
}

func (a App) GetUser(id int64) (user.User, error) {
	u, err := a.userRepository.GetCurrentUser(id)
	if err != nil {
		return u, NoSuchUserError{ErrorNoSuchUser}
	}
	return u, nil
}

func (a App) DeleteUser(id int64) error {
	currentPost, err := a.userRepository.GetCurrentUser(id)
	if err != nil {
		return NoSuchUserError{ErrorNoSuchUser}
	}
	return a.userRepository.DeleteUser(currentPost)
}
