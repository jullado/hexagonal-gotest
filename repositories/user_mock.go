package repositories

import (
	"hexagonal-gotest/models"

	"github.com/stretchr/testify/mock"
)

type userRepoMock struct {
	mock.Mock
}

func NewUserRepoMock() userRepoMock {
	return userRepoMock{}
}

func (m *userRepoMock) Gets(filter models.RepoGetUserModel) (result []models.RepoUserModel, err error) {
	args := m.Called(filter)
	res, ok := args.Get(0).([]models.RepoUserModel)
	if !ok {
		return nil, args.Error(1)
	}
	return res, args.Error(1)
}

func (m *userRepoMock) Create(payload models.RepoCreateUserModel) (err error) {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *userRepoMock) Update(userId string, payload models.RepoUpdateUserModel) (err error) {
	args := m.Called(userId, payload)
	return args.Error(0)
}

func (m *userRepoMock) Delete(userId string) (err error) {
	args := m.Called(userId)
	return args.Error(0)
}
