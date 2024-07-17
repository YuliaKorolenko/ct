package app

import (
	validator "github.com/YuliaKorolenko/valid"
	"homework6/internal/ads"
)

type App struct {
	repository Repository
}

type Repository interface {
	CreateAd(ad ads.Ad) (ads.Ad, error)
	ChangeAdStatus(ad ads.Ad) (ads.Ad, error)
	UpdateAd(ad ads.Ad) (ads.Ad, error)
	GetCurrentAd(id int64) ads.Ad
}

func NewApp(repo Repository) App {
	return App{repository: repo}
}

func (a App) CreateAd(title string, text string, userId int64) (ads.Ad, error) {
	post := ads.New(-1, title, text, userId)
	err := validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.CreateAd(post)
}

func (a App) ChangeAdStatus(id int64, userId int64, isPublished bool) (ads.Ad, error) {
	currentPost, err := a.checkAuthor(id, userId)
	if err != nil {
		return currentPost, err
	}
	post := ads.Update(id, currentPost.Title, currentPost.Text, userId, isPublished)
	err = validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.ChangeAdStatus(post)
}

func (a App) UpdateAd(id int64, userId int64, title string, text string) (ads.Ad, error) {
	currentPost, err := a.checkAuthor(id, userId)
	if err != nil {
		return currentPost, err
	}
	post := ads.Update(id, title, text, userId, currentPost.Published)
	err = validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.UpdateAd(post)
}

func (a App) checkAuthor(id int64, userId int64) (ads.Ad, error) {
	currentPost := a.repository.GetCurrentAd(id)
	if currentPost.AuthorID != userId {
		return currentPost, PermissionDeniedError{Err: ErrorNoRights}
	}
	return currentPost, nil
}
