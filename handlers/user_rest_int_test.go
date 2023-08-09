package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"hexagonal-gotest/handlers"
	"hexagonal-gotest/models"
	"hexagonal-gotest/repositories"
	"hexagonal-gotest/services"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
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
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return([]models.RepoUserModel{
					{Username: tt.body.Username},
				}, nil)
			case "error500":
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return(nil, errors.New(""))
			default:
				userRepo.On("Gets", mock.AnythingOfType("models.RepoGetUserModel")).Return([]models.RepoUserModel{}, nil)
			}

			// mock Create user
			switch tt.name {
			case "error500":
				userRepo.On("Create", mock.AnythingOfType("models.RepoCreateUserModel")).Return(errors.New(""))
			default:
				userRepo.On("Create", mock.AnythingOfType("models.RepoCreateUserModel")).Return(nil)
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

			// -------------------- Assert (ยืนยีน) --------------------

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}
