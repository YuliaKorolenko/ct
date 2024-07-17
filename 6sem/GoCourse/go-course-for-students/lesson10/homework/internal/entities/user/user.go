package user

//id, nickname, email

type User struct {
	ID       int64
	Nickname string `validate:"min:1;max:100"`
	Email    string `validate:"min:1;max:100"`
}

func NewUser(id int64, nickname string, email string) User {
	return User{ID: id, Nickname: nickname, Email: email}
}
