package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
)

type User struct {
	Email          string
	PasswordDigest string
	FavoriteCake   string
}
type UserRepository interface {
	Add(string, User) error
	Get(string) (User, error)
	Update(string, User) error
	Delete(string) (User, error)
}

type UserService struct {
	repository UserRepository
}
type UserRegisterParams struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FavoriteCake string `json:"favorite_cake"`
}

func validateCake(cake string) error {
	// 3. Favorite cake not empty
	// 4. Favorite cake only alphabetic
	if len(cake) < 1 {
		return errors.New("favourite cake can not be empty")
	}
	for _, charVariable := range cake {
		if (charVariable > 'z' || charVariable < 'a') && (charVariable > 'Z' || charVariable < 'A') {
			return errors.New("favourite cake can contain only alphabetic chars")
		}
	}
	return nil
}

func validateEmail(email string) error {
	// 1. Email is valid
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("provide a valid email please")
	}
	return nil
}

func validatePassword(pass string) error {
	// 2. Password at least 8 symbols
	if len(pass) < 8 {
		err := errors.New("password is too short (min 8 symbols)")
		return err
	}
	return nil
}

func validateRegisterParams(p *UserRegisterParams) error {
	if err := validatePassword(p.Password); err != nil {
		return err
	}

	if err := validateEmail(p.Email); err != nil {
		return err
	}

	if err := validateCake(p.FavoriteCake); err != nil {
		return err
	}
	return nil
}
func (u *UserService) Register(w http.ResponseWriter, r *http.Request) {
	params := &UserRegisterParams{}
	err := json.NewDecoder(r.Body).Decode(params)
	if err != nil {
		handleError(errors.New("could not read params"), w)
		return
	}
	if err := validateRegisterParams(params); err != nil {
		handleError(err, w)
		return
	}
	passwordDigest := md5.New().Sum([]byte(params.Password))
	newUser := User{
		Email:          params.Email,
		PasswordDigest: string(passwordDigest),
		FavoriteCake:   params.FavoriteCake,
	}
	err = u.repository.Add(params.Email, newUser)
	if err != nil {
		handleError(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("registered"))
}

func handleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write([]byte(err.Error()))
}
