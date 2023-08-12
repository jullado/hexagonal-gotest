package services

import (
	"errors"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"
	"hexagonal-gotest/utils"

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

	return
}

func (s userSrv) Login(username, password string) (token string, err error) {
	if username == "" {
		return token, errors.New(models.ErrUsernameNotfound)
	}

	if password == "" {
		return token, errors.New(models.ErrPasswordNotfound)
	}

	if len(password) < 6 || len(password) > 16 {
		return token, errors.New(models.ErrPasswordFormat)
	}

	resUsers, err := s.userRepo.Gets(models.RepoGetUserModel{Username: username})
	if err != nil {
		return token, errors.New(models.ErrUnexpected)
	}
	if len(resUsers) == 0 {
		return token, errors.New(models.ErrUsernameIsNotExist)
	}

	if resUsers[0].Password != password {
		return token, errors.New(models.ErrUnauthorized)
	}

	token, err = utils.SignJWT(utils.TokenDataModel{UserId: resUsers[0].UserId, Username: resUsers[0].Username})

	return
}

func (s userSrv) ResetPassword(userId, newPassword string) (err error) {
	if newPassword == "" {
		return errors.New(models.ErrPasswordNotfound)
	}

	if len(newPassword) < 6 || len(newPassword) > 16 {
		return errors.New(models.ErrPasswordFormat)
	}

	if _, err := uuid.Parse(userId); err != nil {
		return errors.New(models.ErrUserIdFormat)
	}

	err = s.userRepo.Update(userId, models.RepoUpdateUserModel{Password: newPassword})
	if err != nil {
		if err.Error() == models.ErrUserIdIsNotExist {
			return errors.New(models.ErrUserIdIsNotExist)
		}

		return errors.New(models.ErrUnexpected)
	}

	return
}

func (s userSrv) DeleteUser(userId string) (err error) {
	if _, err := uuid.Parse(userId); err != nil {
		return errors.New(models.ErrUserIdFormat)
	}

	err = s.userRepo.Delete(userId)
	if err != nil {
		if err.Error() == models.ErrUserIdIsNotExist {
			return errors.New(models.ErrUserIdIsNotExist)
		}

		return errors.New(models.ErrUnexpected)
	}

	return
}
