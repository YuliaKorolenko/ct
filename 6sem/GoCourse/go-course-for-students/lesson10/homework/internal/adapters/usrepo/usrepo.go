package usrepo

import (
	"errors"
	"homework10/internal/app"
	"homework10/internal/entities/user"
	"sync"
)

func NewUserRepo() app.URepository {
	return &SliceRepositoryUser{sliceUser: make(map[int64]user.User), mx: new(sync.RWMutex), currentId: 0}
}

type SliceRepositoryUser struct {
	mx        *sync.RWMutex
	sliceUser map[int64]user.User
	currentId int64
}

func (s *SliceRepositoryUser) CreateUser(u user.User) (user.User, error) {
	newUser := user.NewUser(s.currentId, u.Nickname, u.Email)
	s.currentId++
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
	if s.sliceUser[id] == s.sliceUser[s.currentId+1] {
		return s.sliceUser[id], errors.New("no such id in ad")
	}
	return s.sliceUser[id], nil
}

func (s *SliceRepositoryUser) DeleteUser(user user.User) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.sliceUser, user.ID)
	return nil
}
