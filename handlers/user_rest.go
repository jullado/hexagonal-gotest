package handlers

import (
	"hexagonal-gotest/models"
	"hexagonal-gotest/services"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) userHandler {
	return userHandler{userSrv}
}

func (h userHandler) Register(c *fiber.Ctx) error {
	body := models.HandRegisterBodyModel{}
	c.BodyParser(&body)

	err := h.userSrv.Register(body.Username, body.Password)
	if err != nil {
		switch err.Error() {
		case models.ErrUnexpected:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "register success",
	})
}

func (h userHandler) Login(c *fiber.Ctx) error {
	body := models.HandLoginBodyModel{}
	c.BodyParser(&body)

	token, err := h.userSrv.Login(body.Username, body.Password)
	if err != nil {
		switch err.Error() {
		case models.ErrUsernameIsNotExist, models.ErrUnauthorized:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"token":   token,
				"message": err.Error(),
			})

		case models.ErrUnexpected:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"token":   token,
				"message": err.Error(),
			})

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"token":   token,
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"message": "login success",
	})
}

func (h userHandler) ResetPassword(c *fiber.Ctx) error {
	params := models.HandResetPasswordParamsModel{}
	c.ParamsParser(&params)

	body := models.HandResetPasswordBodyModel{}
	c.BodyParser(&body)

	err := h.userSrv.ResetPassword(params.UserId, body.Password)
	if err != nil {
		switch err.Error() {
		case models.ErrUnexpected:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "reset password success",
	})
}

func (h userHandler) DeleteUser(c *fiber.Ctx) error {
	params := models.HandDeleteUserParamsModel{}
	c.ParamsParser(&params)

	err := h.userSrv.DeleteUser(params.UserId)
	if err != nil {
		switch err.Error() {
		case models.ErrUnexpected:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})

		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "delete user success",
	})
}
