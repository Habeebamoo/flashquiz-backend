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

func (user *User) Validate() error {
	if user.Name == "" {
		return errors.New("name is not defined")
	} else if user.Email == "" {
		return errors.New("email is not defined")
	} else if user.Password == "" {
		return errors.New("password is not defined")
	} else {
		return nil
	}
}

func (user *User) LoginValidate() error {
	if user.Email == "" {
		return errors.New("email is not defined")
	} else if user.Password == "" {
		return errors.New("password is not defined")
	} else {
		return nil
	}
}