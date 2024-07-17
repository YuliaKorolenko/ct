package app

import (
	validator "github.com/YuliaKorolenko/valid"
	ads2 "homework8/internal/entities/ads"
	"homework8/internal/entities/user"
)

type App struct {
	repository     Repository
	userRepository URepository
}

type Repository interface {
	CreateAd(ad ads2.Ad) (ads2.Ad, error)
	ChangeAdStatus(ad ads2.Ad) (ads2.Ad, error)
	UpdateAd(ad ads2.Ad) (ads2.Ad, error)
	GetAdsList(ad ads2.AdFilter) ([]ads2.Ad, error)
	GetCurrentAd(id int64) (ads2.Ad, error)
}

func NewApp(repo Repository, urepo URepository) App {
	return App{repository: repo, userRepository: urepo}
}

func (a App) CreateAd(title string, text string, userId int64) (ads2.Ad, error) {
	post := ads2.New(-1, title, text, userId)
	_, err := a.userRepository.GetCurrentUser(userId)
	if err != nil {
		return post, NoSuchUserError{ErrorNoSuchUser}
	}
	err = validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.CreateAd(post)
}

func (a App) GetAdsList(filter ads2.AdFilter) ([]ads2.Ad, error) {
	return a.repository.GetAdsList(filter)
}

func (a App) GetAd(id int64) (ads2.Ad, error) {
	ad, err := a.repository.GetCurrentAd(id)
	if err != nil {
		return ad, NoSuchAdError{Err: ErrorNoAd}
	}
	return ad, nil
}

func (a App) ChangeAdStatus(id int64, userId int64, isPublished bool) (ads2.Ad, error) {
	currentPost, err := a.checkAuthor(id, userId)
	if err != nil {
		return currentPost, err
	}
	post := ads2.Update(id, currentPost.Title, currentPost.Text, userId, isPublished, currentPost.Time)
	err = validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.ChangeAdStatus(post)
}

func (a App) UpdateAd(id int64, userId int64, title string, text string) (ads2.Ad, error) {
	currentPost, err := a.checkAuthor(id, userId)
	if err != nil {
		return currentPost, err
	}
	post := ads2.Update(id, title, text, userId, currentPost.Published, currentPost.Time)
	err = validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.UpdateAd(post)
}

func (a App) checkAuthor(id int64, userId int64) (ads2.Ad, error) {
	currentPost, err := a.repository.GetCurrentAd(id)
	if err != nil {
		return currentPost, NoSuchAdError{Err: ErrorNoAd}
	}
	_, err = a.userRepository.GetCurrentUser(userId)
	if err != nil {
		return currentPost, NoSuchUserError{ErrorNoSuchUser}
	}

	if currentPost.AuthorID != userId {
		return currentPost, PermissionDeniedError{Err: ErrorNoRights}
	}
	return currentPost, nil
}

type URepository interface {
	CreateUser(user user.User) (user.User, error)
	UpdateUser(user user.User) (user.User, error)
	GetCurrentUser(id int64) (user.User, error)
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
