package controller

import (
	"micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadPhoto(c *fiber.Ctx) error
}

type UploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// @Summary Upload Photo
// @Description Upload a photo
// @Tags User Upload Photo
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Photo"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/upload/photo [post]
// @Security Bearer
func (u *UploadController) UploadPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get file",
		})		
	}

	result, err := u.fileUploadHelper.UploadPhoto(c.Context(), file)
	if err != nil {
		log.Errorf("[UploadController] UploadPhoto - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file",
		})		
	}

	uploadResponse := response.UploadPhotoResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"data":    uploadResponse,
	})
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &UploadController{
		fileUploadHelper: fileUploadHelper,
	}
}
