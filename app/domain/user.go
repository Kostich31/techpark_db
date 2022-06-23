package domain

type User struct {
	Nickname string `json:"nickname,omitempty"`
	FullName string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type UserUpdate struct {
	FullName string `json:"fullname" validate:"required"`
	About    string `json:"about" validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}

type UserRepository interface {
	AddUser(user User) (User, error)
	GetUser(nickname string) (User, error)
	UpdateUser(user User) (User, error)
	GetUsersByNicknameOrEmail(nickname string, email string) ([]User, error)
}

type UserUseCase interface {
	CreateUser(user User) ([]User, error)
	GetUserProfile(nickname string) (User, *CustomError)
	UpdateUserProfile(user User) (User, *CustomError)
}
