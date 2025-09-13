package user

type UserRegisterInput struct {
	Name     string
	Email    string
	Password string
}

type UserLoginInput struct {
	Email    string
	Password string
}
