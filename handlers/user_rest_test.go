package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"hexagonal-gotest/handlers"
	"hexagonal-gotest/models"
	"hexagonal-gotest/services"
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

			// -------------------- Assert (ยืนยีน) --------------------
			userSrv.AssertCalled(t, "Register", tt.body.Username, tt.body.Password)

			assert.Equal(t, tt.wantStatusCode, res.StatusCode)

			b, _ := io.ReadAll(res.Body)
			resBody := responseData{}
			json.Unmarshal(b, &resBody)
			assert.Equal(t, tt.wantData, resBody)
		})
	}
}
