package controller

import (
	"micro-warehouse/notification-service/controller/request"
	"micro-warehouse/notification-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type EmailController struct {
	emailUsecase *usecase.EmailUsecase
}

func NewEmailController(emailUsecase *usecase.EmailUsecase) *EmailController {
	return &EmailController{
		emailUsecase: emailUsecase,
	}
}

func (e *EmailController) SendEmail(c *fiber.Ctx) error {
	var req request.SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[EmailController] SendEmail - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validasi request
	if err := validate.Struct(&req); err != nil {
		log.Errorf("[EmailController] SendEmail - validation: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	if err := e.emailUsecase.SendEmail(c.Context(), req); err != nil {
		log.Errorf("[EmailController] SendEmail - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

func (e *EmailController) SendWelcomeEmail(c *fiber.Ctx) error {
	var req request.SendWelcomeEmailRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[EmailController] SendWelcomeEmail - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validasi request
	if err := validate.Struct(&req); err != nil {
		log.Errorf("[EmailController] SendWelcomeEmail - validation: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	if err := e.emailUsecase.SendWelcomeEmail(c.Context(), req); err != nil {
		log.Errorf("[EmailController] SendWelcomeEmail - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send email",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}
