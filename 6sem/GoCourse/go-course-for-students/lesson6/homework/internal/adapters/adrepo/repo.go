package adrepo

import (
	"homework6/internal/ads"
	"homework6/internal/app"
)

func New() app.Repository {
	return SliceRepository{sliceAd: make(map[int64]ads.Ad)}
}

type SliceRepository struct {
	sliceAd map[int64]ads.Ad
}

func (s SliceRepository) CreateAd(ad ads.Ad) (ads.Ad, error) {
	newAd := ads.New(int64(len(s.sliceAd)), ad.Title, ad.Text, ad.AuthorID)
	s.sliceAd[newAd.ID] = newAd
	return newAd, nil
}

func (s SliceRepository) ChangeAdStatus(ad ads.Ad) (ads.Ad, error) {
	s.sliceAd[ad.ID] = ads.Update(ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published)
	return s.sliceAd[ad.ID], nil
}

func (s SliceRepository) UpdateAd(ad ads.Ad) (ads.Ad, error) {
	s.sliceAd[ad.ID] = ads.Update(ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published)
	return s.sliceAd[ad.ID], nil
}

func (s SliceRepository) GetCurrentAd(id int64) ads.Ad {
	return s.sliceAd[id]
}
