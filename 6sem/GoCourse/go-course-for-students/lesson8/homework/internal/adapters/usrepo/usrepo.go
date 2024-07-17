package usrepo

import (
	"errors"
	"homework8/internal/app"
	"homework8/internal/entities/user"
	"sync"
)

func NewUserRepo() app.URepository {
	return &SliceRepositoryUser{sliceUser: make(map[int64]user.User), mx: new(sync.RWMutex)}
}

type SliceRepositoryUser struct {
	mx        *sync.RWMutex
	sliceUser map[int64]user.User
}

func (s *SliceRepositoryUser) CreateUser(u user.User) (user.User, error) {
	newUser := user.NewUser(int64(len(s.sliceUser)), u.Nickname, u.Email)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceUser[newUser.ID] = newUser
	return s.sliceUser[newUser.ID], nil
}

func (s *SliceRepositoryUser) UpdateUser(u user.User) (user.User, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.sliceUser[u.ID] = user.NewUser(u.ID, u.Nickname, u.Email)
	return s.sliceUser[u.ID], nil
}

func (s *SliceRepositoryUser) GetCurrentUser(id int64) (user.User, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if id >= int64(len(s.sliceUser)) {
		return s.sliceUser[id], errors.New("no such id in ad")
	}
	return s.sliceUser[id], nil
}
