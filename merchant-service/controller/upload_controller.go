package controller

import (
	"micro-warehouse/merchant-service/controller/response"
	"micro-warehouse/merchant-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadMerchantPhoto(c *fiber.Ctx) error
}

type uploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// UploadMerchantPhoto implements UploadControllerInterface.
// @Summary Upload Merchant Photo
// @Tags Merchant Upload Photo
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Photo"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/upload-merchant [post]
// @Security Bearer
func (u *uploadController) UploadMerchantPhoto(c *fiber.Ctx) error {
	file, err := c.FormFile("photo")
	if err != nil {
		log.Errorf("[UploadController] UploadMerchantPhoto - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No file uploaded",
			"error": err.Error(),
		})
	}

	// Upload to Supabse using FileUploadHelper
	result, err := u.fileUploadHelper.UploadPhoto(c.Context(), file)
	if err != nil {
		log.Errorf("[UploadController] UploadMerchantPhoto - 2: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file",
			"error": err.Error(),
		})
	}

	// Create Response
	uploadResponse := response.UploadResponse{
		URL: result.URL,
		Filename: result.Filename,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"data": uploadResponse,
	})
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &uploadController{
		fileUploadHelper: fileUploadHelper,
	}
}
