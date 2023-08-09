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
