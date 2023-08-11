package services

import "github.com/stretchr/testify/mock"

type userSrvMock struct {
	mock.Mock
}

func NewUserSrvMock() userSrvMock {
	return userSrvMock{}
}

func (m *userSrvMock) Register(username, password string) (err error) {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *userSrvMock) Login(username, password string) (token string, err error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *userSrvMock) ResetPassword(userId, newPassword string) (err error) {
	args := m.Called(newPassword)
	return args.Error(0)
}

func (m *userSrvMock) DeleteUser(userId string) (err error) {
	args := m.Called(userId)
	return args.Error(0)
}
