package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hexagonal-gotest/handlers"
	"hexagonal-gotest/models"
	"hexagonal-gotest/services"
	"hexagonal-gotest/utils"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
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
			name: "error400",
			body: reqBody{Username: "", Password: "admin"},
			wantData: responseData{
				Message: models.ErrUsernameNotfound,
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
			userSrv := services.NewUserSrvMock()

			// mock register service
			switch tt.name {
			case "success":
				userSrv.On("Register", tt.body.Username, tt.body.Password).Return(nil)
			default:
				userSrv.On("Register", tt.body.Username, tt.body.Password).Return(errors.New(tt.wantData.Message))
			}

			userHandler := handlers.NewUserHandler(&userSrv)

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
			userSrv.AssertCalled(t, "Register", tt.body.Username, tt.body.Password)

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}

func TestLogin(t *testing.T) {
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type responseData struct {
		Token   string `json:"token"`
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
			name: "error400",
			body: reqBody{Username: "", Password: "admin"},
			wantData: responseData{
				Token:   "",
				Message: models.ErrUsernameNotfound,
			},
			wantStatusCode: 400,
		},
		{
			name: "error500",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Token:   "",
				Message: models.ErrUnexpected,
			},
			wantStatusCode: 500,
		},
		{
			name: "success",
			body: reqBody{Username: "admin", Password: "admin01"},
			wantData: responseData{
				Token: func() string {
					token, _ := utils.SignJWT(utils.TokenDataModel{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191", Username: "admin"})
					return token
				}(),
				Message: "login success",
			},
			wantStatusCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ------------------- Arrange (เตรียมของ) --------------------
			userSrv := services.NewUserSrvMock()

			// mock login service
			switch tt.name {
			case "success":
				userSrv.On("Login", tt.body.Username, tt.body.Password).Return(tt.wantData.Token, nil)
			default:
				userSrv.On("Login", tt.body.Username, tt.body.Password).Return(tt.wantData.Token, errors.New(tt.wantData.Message))
			}

			userHandler := handlers.NewUserHandler(&userSrv)

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
			userSrv.AssertCalled(t, "Login", tt.body.Username, tt.body.Password)

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}

func TestResetPassword(t *testing.T) {
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
			name:   "error400",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			body:   reqBody{Password: "admin"},
			wantData: responseData{
				Message: models.ErrUsernameNotfound,
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
			userSrv := services.NewUserSrvMock()

			// mock reset password service
			switch tt.name {
			case "success":
				userSrv.On("ResetPassword", tt.params.UserId, tt.body.Password).Return(nil)
			default:
				userSrv.On("ResetPassword", tt.params.UserId, tt.body.Password).Return(errors.New(tt.wantData.Message))
			}

			userHandler := handlers.NewUserHandler(&userSrv)

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
			userSrv.AssertCalled(t, "ResetPassword", tt.params.UserId, tt.body.Password)

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}

func TestDeleteUser(t *testing.T) {
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
			name:   "error400",
			params: reqParams{UserId: "225cfc88-c66b-4f2f-b424-a3b74e3b1191"},
			wantData: responseData{
				Message: models.ErrUsernameNotfound,
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
			userSrv := services.NewUserSrvMock()

			// mock reset password service
			switch tt.name {
			case "success":
				userSrv.On("DeleteUser", tt.params.UserId).Return(nil)
			default:
				userSrv.On("DeleteUser", tt.params.UserId).Return(errors.New(tt.wantData.Message))
			}

			userHandler := handlers.NewUserHandler(&userSrv)

			// http request
			app := fiber.New()
			app.Delete("/delete/:user_id", userHandler.DeleteUser)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/delete/%v", tt.params.UserId), nil)

			// -------------------- Act (กระทำ)--------------------
			res, _ := app.Test(req)
			defer res.Body.Close()

			// -------------------- Assert (ยืนยัน) --------------------
			userSrv.AssertCalled(t, "DeleteUser", tt.params.UserId)

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}
