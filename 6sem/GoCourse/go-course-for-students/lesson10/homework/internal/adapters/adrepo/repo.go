package adrepo

import (
	"errors"
	"homework10/internal/app"
	ads "homework10/internal/entities/ads"
	"strings"
	"sync"
	"time"
)

func New() app.Repository {
	return &SliceRepository{sliceAd: make(map[int64]ads.Ad), mx: new(sync.RWMutex), currentId: 0}
}

type SliceRepository struct {
	mx        *sync.RWMutex
	sliceAd   map[int64]ads.Ad
	currentId int64
}

func (s *SliceRepository) CreateAd(ad ads.Ad) (ads.Ad, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	newAd := ads.New(s.currentId, ad.Title, ad.Text, ad.AuthorID)
	s.currentId++
	s.sliceAd[newAd.ID] = newAd
	return newAd, nil
}

func (s *SliceRepository) GetAdsListSlow(f ads.AdFilter) ([]ads.Ad, error) {
	answer := make([]ads.Ad, 0)
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, val := range s.sliceAd {
		maxLen := len(val.Title)
		if maxLen > len(f.Title) {
			maxLen = len(f.Title)
		}
		isEqual := false
		for i := 0; i < maxLen; i++ {
			if len(val.Title) > i && len(f.Title) > i {
				if val.Title[i] == f.Title[i] && i == 0 {
					isEqual = true
				}
			}
		}
		if val.Published == f.Published {
			if val.AuthorID == f.AuthorID || f.AuthorID == -1 {
				if val.Time.Format(time.DateOnly) == f.Time || f.Time == "" {
					if isEqual || len(f.Title) == 0 {
						answer = append(answer, val)
					}
				}
			}
		}

	}
	if len(answer) == 0 {
		return nil, nil
	}
	return answer, nil
}

func (s *SliceRepository) GetAdsList(f ads.AdFilter) ([]ads.Ad, error) {
	answer := make([]ads.Ad, 0)
	s.mx.RLock()
	defer s.mx.RUnlock()
	for _, val := range s.sliceAd {
		if (strings.HasPrefix(val.Title, f.Title)) &&
			(val.Published == f.Published) &&
			(val.AuthorID == f.AuthorID || f.AuthorID == -1) &&
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

func (s *SliceRepository) ChangeAdStatus(ad ads.Ad) (ads.Ad, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceAd[ad.ID] = ads.Update(ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published, ad.Time)
	return s.sliceAd[ad.ID], nil
}

func (s *SliceRepository) UpdateAd(ad ads.Ad) (ads.Ad, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceAd[ad.ID] = ads.Update(ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published, ad.Time)
	return s.sliceAd[ad.ID], nil
}

func (s *SliceRepository) GetCurrentAd(id int64) (ads.Ad, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if s.sliceAd[id] == s.sliceAd[s.currentId+1] {
		return s.sliceAd[id], errors.New("no such id in ad")
	}
	return s.sliceAd[id], nil
}

func (s *SliceRepository) DeleteAd(ad ads.Ad) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.sliceAd, ad.ID)
	return nil
}
