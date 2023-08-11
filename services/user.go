package services

// PORT user service
type UserService interface {
	Register(username, password string) (err error)

	Login(username, password string) (token string, err error)

	ResetPassword(userId, newPasword string) (err error)

	DeleteUser(userId string) (err error)
}
