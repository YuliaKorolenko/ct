package app

import (
	validator "github.com/YuliaKorolenko/valid"
	"homework10/internal/entities/ads"
)

//go:generate mockery --name Repository
type Repository interface {
	CreateAd(ad ads.Ad) (ads.Ad, error)
	ChangeAdStatus(ad ads.Ad) (ads.Ad, error)
	UpdateAd(ad ads.Ad) (ads.Ad, error)
	GetAdsList(ad ads.AdFilter) ([]ads.Ad, error)
	GetAdsListSlow(ad ads.AdFilter) ([]ads.Ad, error)
	GetCurrentAd(id int64) (ads.Ad, error)
	DeleteAd(ad ads.Ad) error
}

func NewApp(repo Repository, urepo URepository) App {
	return App{repository: repo, userRepository: urepo}
}

func (a App) CreateAd(title string, text string, userId int64) (ads.Ad, error) {
	post := ads.New(-1, title, text, userId)
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

func (a App) GetAdsList(filter ads.AdFilter) ([]ads.Ad, error) {
	return a.repository.GetAdsList(filter)
}

func (a App) GetAd(id int64) (ads.Ad, error) {
	ad, err := a.repository.GetCurrentAd(id)
	if err != nil {
		return ad, NoSuchAdError{Err: ErrorNoAd}
	}
	_, nil := a.userRepository.GetCurrentUser(ad.AuthorID)
	if err != nil {
		return ad, NoSuchUserError{ErrorNoSuchUser}
	}
	return ad, nil
}

func (a App) ChangeAdStatus(id int64, userId int64, isPublished bool) (ads.Ad, error) {
	currentPost, err := a.checkAuthor(id, userId)
	if err != nil {
		return currentPost, err
	}
	post := ads.Update(id, currentPost.Title, currentPost.Text, userId, isPublished, currentPost.Time)
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
	post := ads.Update(id, title, text, userId, currentPost.Published, currentPost.Time)
	err = validator.Validate(post)
	if err != nil {
		return post, ValidationApError{ErrorValidation}
	}
	return a.repository.UpdateAd(post)
}

func (a App) DeleteAd(id int64, userId int64) error {
	currentPost, err := a.checkAuthor(id, userId)
	if err != nil {
		return err
	}
	return a.repository.DeleteAd(currentPost)
}

func (a App) checkAuthor(id int64, userId int64) (ads.Ad, error) {
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
