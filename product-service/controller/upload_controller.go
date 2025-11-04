package controller

import (
	"micro-warehouse/product-service/controller/response"
	"micro-warehouse/product-service/pkg/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UploadControllerInterface interface {
	UploadProductImage(ctx *fiber.Ctx) error
	UploadCategoryImage(ctx *fiber.Ctx) error
}

type uploadController struct {
	fileUploadHelper *storage.FileUploadHelper
}

// UploadCategoryImage implements UploadControllerInterface.
func (u *uploadController) UploadCategoryImage(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Errorf("failed to get file: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get file",
			"error":   err.Error(),
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx.Context(), file, "categories")
	if err != nil {
		log.Errorf("failed to upload file: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file",
			"error":   err.Error(),
		})
	}

	response := response.UploadResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"data":    response,
	})
}

// UploadProductImage implements UploadControllerInterface.
func (u *uploadController) UploadProductImage(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("image")
	if err != nil {
		log.Errorf("failed to get file category: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to get file category",
			"error":   err.Error(),
		})
	}

	result, err := u.fileUploadHelper.UploadPhoto(ctx.Context(), file, "products")
	if err != nil {
		log.Errorf("failed to upload file product: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to upload file product",
			"error":   err.Error(),
		})
	}

	response := response.UploadResponse{
		URL:      result.URL,
		Path:     result.Path,
		Filename: result.Filename,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"data":    response,
	})
}

func NewUploadController(fileUploadHelper *storage.FileUploadHelper) UploadControllerInterface {
	return &uploadController{
		fileUploadHelper: fileUploadHelper,
	}
}
