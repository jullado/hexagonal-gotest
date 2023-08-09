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
