package controller

import (
	"micro-warehouse/warehouse-service/controller/response"
	"micro-warehouse/warehouse-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadPhoto(c *fiber.Ctx) error
}

type uploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// UploadPhoto implements UploadControllerInterface.
func (u *uploadController) UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get image from request",
			"error":   err.Error(),
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(c.Context(), file)
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
	}

	resp := response.UploadResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file upload successfully",
		"data":    resp,
	})
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &uploadController{
		fileUploadHelper: fileUploadHelper,
	}
}
