package main

import (
	"fmt"
	"sync"

	"github.com/go-chi/jwtauth/v5"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"`
}

type UserData map[string]*User

type Users struct {
	Data UserData
	Mu   sync.RWMutex
}

type UserStorage interface {
	Create(*User) (*User, error)
	checkUser(*User, *jwtauth.JWTAuth) (*User, error)
	getByLogin(login string) (*User, bool)
}

func newUserStorage() *Users {
	return &Users{
		Data: make(UserData),
	}
}

func (u *Users) Create(user *User) (*User, error) {
	u.Mu.RLock()
	_, ok := u.Data[user.Login]
	u.Mu.RUnlock()

	if ok || user.Login == "" || user.Password == "" {
		return nil, fmt.Errorf("SIGNUP ERROR\n")
	}

	v := &User{Login: user.Login, Password: user.Password}
	u.Mu.Lock()
	u.Data[user.Login] = v
	u.Mu.Unlock()

	return v, nil
}

func (u *Users) checkUser(user *User, token *jwtauth.JWTAuth) (*User, error) {
	u.Mu.RLock()
	v, ok := u.Data[user.Login]
	u.Mu.RUnlock()

	if !ok || v.Password != user.Password {
		return nil, fmt.Errorf("AUTHORIZATION ERROR\n")
	}

	u.Mu.Lock()
	_, v.Token, _ = token.Encode(map[string]interface{}{"login": user.Login})
	u.Mu.Unlock()

	return v, nil
}

func (u *Users) getByLogin(login string) (*User, bool) {
	u.Mu.RLock()
	v, ok := u.Data[login]
	u.Mu.RUnlock()

	return v, ok
}
