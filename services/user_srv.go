package services

import (
	"errors"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"

	"github.com/google/uuid"
)

type userSrv struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return userSrv{userRepo}
}

func (s userSrv) Register(username, password string) (err error) {
	if username == "" {
		return errors.New(models.ErrUsernameNotfound)
	}

	if password == "" {
		return errors.New(models.ErrPasswordNotfound)
	}

	if len(password) < 6 || len(password) > 16 {
		return errors.New(models.ErrPasswordFormat)
	}

	resGets, err := s.userRepo.Gets(models.RepoGetUserModel{Username: username})
	if err != nil {
		return errors.New(models.ErrUnexpected)
	}
	if len(resGets) > 0 {
		return errors.New(models.ErrUsernameIsExist)
	}

	err = s.userRepo.Create(models.RepoCreateUserModel{UserId: uuid.NewString(), Username: username, Password: password})
	if err != nil {
		return errors.New(models.ErrUnexpected)
	}

	return nil
}
