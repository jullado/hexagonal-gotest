package services_test

import (
	"errors"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"
	"hexagonal-gotest/services"
	"hexagonal-gotest/utils"
	"testing"
	"time"

	"github.com/google/uuid"
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
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.args.username
				})).Return([]models.RepoUserModel{
					{Username: tt.args.username},
				}, nil)
			case "unexpected gets user":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return(nil, errors.New(""))
			default:
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.args.username
				})).Return([]models.RepoUserModel{}, nil)
			}

			// mock Create user
			switch tt.name {
			case "unexpected create user":
				userRepo.On("Create", mock.AnythingOfType("models.RepoCreateUserModel")).Return(errors.New(""))
			default:
				userRepo.On("Create", mock.MatchedBy(func(payload models.RepoCreateUserModel) bool {
					_, errUUID := uuid.Parse(payload.UserId)
					return payload.Username == tt.args.username && payload.Password == tt.args.password && errUUID == nil
				})).Return(nil)
			}

			userSrv := services.NewUserService(&userRepo)

			// -------------------- Act (กระทำ)--------------------
			err := userSrv.Register(tt.args.username, tt.args.password)

			// -------------------- Assert (ยืนยัน) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				userRepo.AssertCalled(t, "Create", mock.MatchedBy(func(payload models.RepoCreateUserModel) bool {
					_, errUUID := uuid.Parse(payload.UserId)
					return payload.Username == tt.args.username && payload.Password == tt.args.password && errUUID == nil
				}))
			}
		})
	}
}

func TestLogin(t *testing.T) {
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
			wantErr: errors.New(models.ErrUsernameIsNotExist),
		},
		{
			name:    "error6",
			args:    args{username: "admin", password: "admin01"},
			wantErr: errors.New(models.ErrUnauthorized),
		},
		{
			name:    "unexpected gets user",
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

			// mock Get user
			switch tt.name {
			case "error5":
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.args.username
				})).Return([]models.RepoUserModel{}, nil)

			case "error6":
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.args.username
				})).Return([]models.RepoUserModel{
					{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", Username: tt.args.username, Password: ""},
				}, nil)

			case "unexpected gets user":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return(nil, errors.New(""))

			default:
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.args.username
				})).Return([]models.RepoUserModel{
					{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", Username: tt.args.username, Password: tt.args.password},
				}, nil)
			}

			userService := services.NewUserService(&userRepo)

			// -------------------- Act (กระทำ)--------------------
			gotToken, err := userService.Login(tt.args.username, tt.args.password)

			// -------------------- Assert (ยืนยัน) --------------------
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				userRepo.AssertCalled(t, "Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.args.username
				}))

				result, err := utils.ValidateJWT(gotToken)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, tt.args.username, result.Username)
				assert.LessOrEqual(t, time.Now(), result.ExpiresAt.Time)
				_, err = uuid.Parse(result.UserId)
				assert.NoError(t, err)
			}
		})
	}
}

func TestResetPassword(t *testing.T) {
	type args struct {
		userId      string
		newPassword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "error1",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", newPassword: ""},
			wantErr: errors.New(models.ErrPasswordNotfound),
		},
		{
			name:    "error2",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", newPassword: "123"},
			wantErr: errors.New(models.ErrPasswordFormat),
		},
		{
			name:    "error3",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", newPassword: "123456789123456789"},
			wantErr: errors.New(models.ErrPasswordFormat),
		},
		{
			name:    "error4",
			args:    args{userId: "", newPassword: "admin01"},
			wantErr: errors.New(models.ErrUserIdFormat),
		},
		{
			name:    "error5",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", newPassword: "admin01"},
			wantErr: errors.New(models.ErrUserIdIsNotExist),
		},
		{
			name:    "unexpected update user",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", newPassword: "admin01"},
			wantErr: errors.New(models.ErrUnexpected),
		},
		{
			name:    "success1",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", newPassword: "admin01"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()

			// mock Update user
			switch tt.name {
			case "error5":
				userRepo.On("Update", tt.args.userId, mock.MatchedBy(func(payload models.RepoUpdateUserModel) bool {
					return payload.Password == tt.args.newPassword
				})).Return(errors.New(tt.wantErr.Error()))

			case "unexpected update user":
				userRepo.On("Update", tt.args.userId, mock.AnythingOfType("models.RepoUpdateUserModel")).Return(errors.New(""))

			default:
				userRepo.On("Update", tt.args.userId, mock.MatchedBy(func(payload models.RepoUpdateUserModel) bool {
					return payload.Password == tt.args.newPassword
				})).Return(nil)
			}

			userService := services.NewUserService(&userRepo)

			// -------------------- Act (กระทำ)--------------------
			err := userService.ResetPassword(tt.args.userId, tt.args.newPassword)

			// -------------------- Assert (ยืนยัน) --------------------
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				userRepo.AssertCalled(t, "Update", tt.args.userId, mock.MatchedBy(func(filter models.RepoUpdateUserModel) bool {
					return filter.Password == tt.args.newPassword
				}))
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		userId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name:    "error1",
			args:    args{userId: ""},
			wantErr: errors.New(models.ErrUserIdFormat),
		},
		{
			name:    "error2",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantErr: errors.New(models.ErrUserIdIsNotExist),
		},
		{
			name:    "unexpected delete user",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantErr: errors.New(models.ErrUnexpected),
		},
		{
			name:    "success1",
			args:    args{userId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()

			// mock Delete user
			switch tt.name {
			case "error2":
				userRepo.On("Delete", tt.args.userId).Return(errors.New(tt.wantErr.Error()))

			case "unexpected delete user":
				userRepo.On("Delete", tt.args.userId).Return(errors.New(""))

			default:
				userRepo.On("Delete", tt.args.userId).Return(nil)
			}

			userService := services.NewUserService(&userRepo)

			// -------------------- Act (กระทำ)--------------------
			err := userService.DeleteUser(tt.args.userId)

			// -------------------- Assert (ยืนยัน) --------------------
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				userRepo.AssertCalled(t, "Delete", tt.args.userId)
			}
		})
	}
}
