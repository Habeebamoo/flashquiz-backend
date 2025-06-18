package models

import "errors"

type User struct {
	Id  int  `json:"id"`
	UserId  string  `json:"userId"`
	Name  string  `json:"name"`
	Email  string  `json:"email"`
	Password  string  `json:"password"`
	IsVerified  bool  `json:"isVerified"`
}

type UserResponse struct {
	UserId  string  `json:"userId"`
	Name  string  `json:"name"`
	Email  string  `json:"email"`
	IsVerified  bool  `json:"isVerified"`
}

type ResetForm struct {
	Password  string  `json:"password"`
}

func (user *User) Validate() error {
	if user.Name == "" {
		return errors.New("name is not missing")
	} else if user.Email == "" {
		return errors.New("email is missing")
	} else if user.Password == "" {
		return errors.New("password is missing")
	} else {
		return nil
	}
}

func (user *User) LoginValidate() error {
	if user.Email == "" {
		return errors.New("email is missing")
	} else if user.Password == "" {
		return errors.New("password is missing")
	} else {
		return nil
	}
}

func (form *ResetForm) Validate() error {
	if form.Password == "" {
		return errors.New("password is missing")
	}
	return nil
}