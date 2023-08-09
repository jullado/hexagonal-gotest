package services

// PORT user service
type UserService interface {
	Register(username, password string) (err error)
}
