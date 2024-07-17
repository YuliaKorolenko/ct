package adrepo

import (
	"errors"
	"homework8/internal/app"
	ads2 "homework8/internal/entities/ads"
	"strings"
	"sync"
	"time"
)

func New() app.Repository {
	return &SliceRepository{sliceAd: make(map[int64]ads2.Ad), mx: new(sync.RWMutex)}
}

type SliceRepository struct {
	mx      *sync.RWMutex
	sliceAd map[int64]ads2.Ad
}

func (s *SliceRepository) CreateAd(ad ads2.Ad) (ads2.Ad, error) {
	newAd := ads2.New(int64(len(s.sliceAd)), ad.Title, ad.Text, ad.AuthorID)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceAd[newAd.ID] = newAd
	return newAd, nil
}

func (s *SliceRepository) GetAdsList(f ads2.AdFilter) ([]ads2.Ad, error) {
	answer := make([]ads2.Ad, 0)
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, val := range s.sliceAd {
		if (val.Published == f.Published) &&
			(val.AuthorID == f.AuthorID || f.AuthorID == -1) &&
			(strings.HasPrefix(val.Title, f.Title)) &&
			(val.Time.Format(time.DateOnly) == f.Time || f.Time == "") {
			{
				answer = append(answer, val)
			}
		}
	}
	if len(answer) == 0 {
		return nil, nil
	}
	return answer, nil
}

func (s *SliceRepository) ChangeAdStatus(ad ads2.Ad) (ads2.Ad, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceAd[ad.ID] = ads2.Update(ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published, ad.Time)
	return s.sliceAd[ad.ID], nil
}

func (s *SliceRepository) UpdateAd(ad ads2.Ad) (ads2.Ad, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceAd[ad.ID] = ads2.Update(ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published, ad.Time)
	return s.sliceAd[ad.ID], nil
}

func (s *SliceRepository) GetCurrentAd(id int64) (ads2.Ad, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if id >= int64(len(s.sliceAd)) {
		return s.sliceAd[id], errors.New("no such id in ad")
	}
	return s.sliceAd[id], nil
}
