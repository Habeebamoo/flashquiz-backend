package user

type User struct {
	Id  int  `json:"user_id"`
	Name  string  `json:"name"`
	Email  string  `json:"email"`
	Password  string  `json:"password"`
}