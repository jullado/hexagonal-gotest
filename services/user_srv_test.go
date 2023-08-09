package services_test

import (
	"errors"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"
	"hexagonal-gotest/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "error1",
			args:    args{username: "", password: "admin01"},
			wantErr: errors.New(models.ErrUsernameNotfound),
		},
		{
			name:    "error2",
			args:    args{username: "admin", password: ""},
			wantErr: errors.New(models.ErrPasswordNotfound),
		},
		{
			name:    "error3",
			args:    args{username: "admin", password: "123"},
			wantErr: errors.New(models.ErrPasswordFormat),
		},
		{
			name:    "error4",
			args:    args{username: "admin", password: "123456789123456789"},
			wantErr: errors.New(models.ErrPasswordFormat),
		},
		{
			name:    "error5",
			args:    args{username: "admin", password: "admin01"},
			wantErr: errors.New(models.ErrUsernameIsExist),
		},
		{
			name:    "unexpected gets user",
			args:    args{username: "admin", password: "admin01"},
			wantErr: errors.New(models.ErrUnexpected),
		},
		{
			name:    "unexpected create user",
			args:    args{username: "admin", password: "admin01"},
			wantErr: errors.New(models.ErrUnexpected),
		},
		{
			name:    "success1",
			args:    args{username: "admin", password: "admin01"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()

			// mock Gets user
			switch tt.name {
			case "error5":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return([]models.RepoUserModel{
					{Username: tt.args.username},
				}, nil)
			case "unexpected gets user":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return(nil, errors.New(""))
			default:
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return([]models.RepoUserModel{}, nil)
			}

			// mock Create user
			switch tt.name {
			case "unexpected create user":
				userRepo.On("Create", mock.AnythingOfType("models.RepoCreateUserModel")).Return(errors.New(""))
			default:
				userRepo.On("Create", mock.AnythingOfType("models.RepoCreateUserModel")).Return(nil)
			}

			userSrv := services.NewUserService(&userRepo)

			// -------------------- Act (กระทำ)--------------------
			err := userSrv.Register(tt.args.username, tt.args.password)

			// -------------------- Assert (ยืนยีน) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				userRepo.AssertCalled(t, "Create", mock.MatchedBy(func(payload models.RepoCreateUserModel) bool {
					return payload.Username == tt.args.username && payload.Password == tt.args.password && payload.UserId != ""
				}))
			}
		})
	}
}
