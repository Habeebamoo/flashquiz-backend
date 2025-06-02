package user

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