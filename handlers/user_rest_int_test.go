package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hexagonal-gotest/handlers"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"
	"hexagonal-gotest/services"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterIntegration(t *testing.T) {
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type responseData struct {
		Message string `json:"message"`
	}
	tests := []struct {
		name           string
		body           reqBody
		wantData       responseData
		wantStatusCode int
	}{
		// TODO: Add test cases.
		{
			name: "error1",
			body: reqBody{Username: "", Password: "admin"},
			wantData: responseData{
				Message: models.ErrUsernameNotfound,
			},
			wantStatusCode: 400,
		},
		{
			name: "error2",
			body: reqBody{Username: "admin", Password: ""},
			wantData: responseData{
				Message: models.ErrPasswordNotfound,
			},
			wantStatusCode: 400,
		},
		{
			name: "error3",
			body: reqBody{Username: "admin", Password: "123"},
			wantData: responseData{
				Message: models.ErrPasswordFormat,
			},
			wantStatusCode: 400,
		},
		{
			name: "error4",
			body: reqBody{Username: "admin", Password: "123456789123456789"},
			wantData: responseData{
				Message: models.ErrPasswordFormat,
			},
			wantStatusCode: 400,
		},
		{
			name: "error5",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUsernameIsExist,
			},
			wantStatusCode: 400,
		},
		{
			name: "error500",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUnexpected,
			},
			wantStatusCode: 500,
		},
		{
			name: "success",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: "register success",
			},
			wantStatusCode: 201,
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
					return filter.Username == tt.body.Username
				})).Return([]models.RepoUserModel{
					{Username: tt.body.Username},
				}, nil)
			case "error500":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return(nil, errors.New(""))
			default:
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.body.Username
				})).Return([]models.RepoUserModel{}, nil)
			}

			// mock Create user
			switch tt.name {
			case "error500":
				userRepo.On("Create", mock.AnythingOfType("models.RepoCreateUserModel")).Return(errors.New(""))
			default:
				userRepo.On("Create", mock.MatchedBy(func(payload models.RepoCreateUserModel) bool {
					_, errUUID := uuid.Parse(payload.UserId)
					return payload.Username == tt.body.Username && payload.Password == tt.body.Password && errUUID == nil
				})).Return(nil)
			}

			userSrv := services.NewUserService(&userRepo)

			userHandler := handlers.NewUserHandler(userSrv)

			// http request
			app := fiber.New()
			app.Post("/register", userHandler.Register)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Add("Content-Type", "application/json")

			// -------------------- Act (กระทำ)--------------------
			res, _ := app.Test(req)
			defer res.Body.Close()

			// -------------------- Assert (ยืนยัน) --------------------

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}

func TestLoginIntegration(t *testing.T) {
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type responseData struct {
		Message string `json:"message"`
	}
	tests := []struct {
		name           string
		body           reqBody
		wantData       responseData
		wantStatusCode int
	}{
		// TODO: Add test cases.
		{
			name: "error1",
			body: reqBody{Username: "", Password: "admin"},
			wantData: responseData{
				Message: models.ErrUsernameNotfound,
			},
			wantStatusCode: 400,
		},
		{
			name: "error2",
			body: reqBody{Username: "admin", Password: ""},
			wantData: responseData{
				Message: models.ErrPasswordNotfound,
			},
			wantStatusCode: 400,
		},
		{
			name: "error3",
			body: reqBody{Username: "admin", Password: "123"},
			wantData: responseData{
				Message: models.ErrPasswordFormat,
			},
			wantStatusCode: 400,
		},
		{
			name: "error4",
			body: reqBody{Username: "admin", Password: "123456789123456789"},
			wantData: responseData{
				Message: models.ErrPasswordFormat,
			},
			wantStatusCode: 400,
		},
		{
			name: "error5",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUsernameIsNotExist,
			},
			wantStatusCode: 401,
		},
		{
			name: "error6",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUnauthorized,
			},
			wantStatusCode: 401,
		},
		{
			name: "error500",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUnexpected,
			},
			wantStatusCode: 500,
		},
		{
			name: "success",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Message: "login success",
			},
			wantStatusCode: 200,
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
					return filter.Username == tt.body.Username
				})).Return([]models.RepoUserModel{}, nil)

			case "error6":
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.body.Username
				})).Return([]models.RepoUserModel{
					{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", Username: tt.body.Username, Password: ""},
				}, nil)

			case "error500":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return(nil, errors.New(""))

			default:
				userRepo.On("Gets", mock.MatchedBy(func(filter models.RepoGetUserModel) bool {
					return filter.Username == tt.body.Username
				})).Return([]models.RepoUserModel{
					{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", Username: tt.body.Username, Password: tt.body.Password},
				}, nil)
			}

			userSrv := services.NewUserService(&userRepo)

			userHandler := handlers.NewUserHandler(userSrv)

			// http request
			app := fiber.New()
			app.Post("/login", userHandler.Login)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Add("Content-Type", "application/json")

			// -------------------- Act (กระทำ)--------------------
			res, _ := app.Test(req)
			defer res.Body.Close()

			// -------------------- Assert (ยืนยัน) --------------------

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}

func TestResetPasswordIntegration(t *testing.T) {
	type reqParams struct {
		UserId string `params:"user_id"`
	}
	type reqBody struct {
		Password string `json:"password"`
	}
	type responseData struct {
		Message string `json:"message"`
	}
	tests := []struct {
		name           string
		params         reqParams
		body           reqBody
		wantData       responseData
		wantStatusCode int
	}{
		// TODO: Add test cases.
		{
			name:   "error1",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: ""},
			wantData: responseData{
				Message: models.ErrPasswordNotfound,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error2",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: "123"},
			wantData: responseData{
				Message: models.ErrPasswordFormat,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error3",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: "123456789123456789"},
			wantData: responseData{
				Message: models.ErrPasswordFormat,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error4",
			params: reqParams{UserId: "123"},
			body:   reqBody{Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUserIdFormat,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error5",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUserIdIsNotExist,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error500",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: "admin01"},
			wantData: responseData{
				Message: models.ErrUnexpected,
			},
			wantStatusCode: 500,
		},
		{
			name:   "success",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: "admin01"},
			wantData: responseData{
				Message: "reset password success",
			},
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()

			// mock Update user
			switch tt.name {
			case "error5":
				userRepo.On("Update", tt.params.UserId, mock.MatchedBy(func(payload models.RepoUpdateUserModel) bool {
					return payload.Password == tt.body.Password
				})).Return(errors.New(tt.wantData.Message))

			case "error500":
				userRepo.On("Update", tt.params.UserId, mock.AnythingOfType("models.RepoUpdateUserModel")).Return(errors.New(""))

			default:
				userRepo.On("Update", tt.params.UserId, mock.MatchedBy(func(payload models.RepoUpdateUserModel) bool {
					return payload.Password == tt.body.Password
				})).Return(nil)
			}

			userSrv := services.NewUserService(&userRepo)

			userHandler := handlers.NewUserHandler(userSrv)

			// http request
			app := fiber.New()
			app.Put("/resetpassword/:user_id", userHandler.ResetPassword)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/resetpassword/%v", tt.params.UserId), bytes.NewBuffer(body))
			req.Header.Add("Content-Type", "application/json")

			// -------------------- Act (กระทำ)--------------------
			res, _ := app.Test(req)
			defer res.Body.Close()

			// -------------------- Assert (ยืนยัน) --------------------
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}

func TestDeleteUserIntegration(t *testing.T) {
	type reqParams struct {
		UserId string `params:"user_id"`
	}
	type responseData struct {
		Message string `json:"message"`
	}
	tests := []struct {
		name           string
		params         reqParams
		wantData       responseData
		wantStatusCode int
	}{
		// TODO: Add test cases.
		{
			name:   "error1",
			params: reqParams{UserId: "123"},
			wantData: responseData{
				Message: models.ErrUserIdFormat,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error2",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantData: responseData{
				Message: models.ErrUserIdIsNotExist,
			},
			wantStatusCode: 400,
		},
		{
			name:   "error500",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantData: responseData{
				Message: models.ErrUnexpected,
			},
			wantStatusCode: 500,
		},
		{
			name:   "success",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantData: responseData{
				Message: "delete user success",
			},
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userRepo := repositories.NewUserRepoMock()

			// mock Delete user
			switch tt.name {
			case "error2":
				userRepo.On("Delete", tt.params.UserId).Return(errors.New(tt.wantData.Message))

			case "error500":
				userRepo.On("Delete", tt.params.UserId).Return(errors.New(""))

			default:
				userRepo.On("Delete", tt.params.UserId).Return(nil)
			}

			userSrv := services.NewUserService(&userRepo)

			userHandler := handlers.NewUserHandler(userSrv)

			// http request
			app := fiber.New()
			app.Delete("/delete/:user_id", userHandler.DeleteUser)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/delete/%v", tt.params.UserId), nil)

			// -------------------- Act (กระทำ)--------------------
			res, _ := app.Test(req)
			defer res.Body.Close()

			// -------------------- Assert (ยืนยัน) --------------------
			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}
