package controller

import (
	"micro-warehouse/product-service/controller/request"
	"micro-warehouse/product-service/controller/response"
	"micro-warehouse/product-service/model"
	"micro-warehouse/product-service/pkg/conv"
	"micro-warehouse/product-service/pkg/pagination"
	"micro-warehouse/product-service/pkg/validator"
	"micro-warehouse/product-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CategoryControllerInterface interface {
	CreateCategory(ctx *fiber.Ctx) error
	GetAllCategories(ctx *fiber.Ctx) error
	GetCategoryByID(ctx *fiber.Ctx) error
	UpdateCategory(ctx *fiber.Ctx) error
	DeleteCategory(ctx *fiber.Ctx) error
}

type categoryController struct {
	categoryUsecase usecase.CategoryUsecaseInterface
}

// CreateCategory implements CategoryControllerInterface.
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param request body request.CreateCategoryRequest true "Category data"
// @Success      201  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/categories [post]
func (c *categoryController) CreateCategory(ctx *fiber.Ctx) error {
	var req request.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[CategoryController] CreateCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[CategoryController] CreateCategory - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Category{
		Name: req.Name,
		Tagline: req.Tagline,
		Photo: req.Photo,
	}

	if err := c.categoryUsecase.CreateCategory(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[CategoryController] CreateCategory - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create category",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Category created successfully",
	})
}

// DeleteCategory implements CategoryControllerInterface.
// @Summary Delete a category by ID
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/categories/{id} [delete]
func (c *categoryController) DeleteCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idUint := conv.StringToUint(id)

	if err := c.categoryUsecase.DeleteCategory(ctx.Context(), idUint); err != nil {
		log.Errorf("[CategoryController] DeleteCategory - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete category",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// GetAllCategories implements CategoryControllerInterface.
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param search query string false "Search keyword"
// @Param sort_by query string false "Sort By"
// @Param sort_order query string false "Sort Order"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/categories [get]
func (c *categoryController) GetAllCategories(ctx *fiber.Ctx) error {
	var req request.GetAllCategoryRequest
	if err := ctx.QueryParser(&req); err != nil {
		log.Errorf("[CategoryController] GetAllCategories - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	categories, total, err := c.categoryUsecase.GetAllCategories(ctx.Context(), req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[CategoryController] GetAllCategories - 2: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get categories",
		})
	}

	pagination := pagination.CalculatePagination(req.Page, req.Limit, int(total))

	var categoriesResponse []response.CategoryResponse
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, response.CategoryResponse{
			ID: category.ID,
			Name: category.Name,
			Tagline: category.Tagline,
			Photo: category.Photo,
		})
	}

	response := response.GetAllCategoryResponse{
		Categories: categoriesResponse,
		Pagination: pagination,
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Categories fetched successfully",
		"data": response,
	})
}

// GetCategoryByID implements CategoryControllerInterface.
// @Summary Get a category by ID
// @Description Get a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/categories/{id} [get]
func (c *categoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idUint := conv.StringToUint(id)

	category, err := c.categoryUsecase.GetCategoryByID(ctx.Context(), idUint)
	if err != nil {
		log.Errorf("[CategoryController] GetCategoryByID - 1: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get category",
		})
	}

	response := response.CategoryResponse{
		ID: category.ID,
		Name: category.Name,
		Tagline: category.Tagline,
		Photo: category.Photo,
		CountProduct: len(category.Products),
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category fetched successfully",
		"data": response,
	})
}

// UpdateCategory implements CategoryControllerInterface.
// @Summary Update a category by ID
// @Description Update a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body request.CreateCategoryRequest true "Category data"
// @Success      200  {object} map[string]interface{}
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /api/v1/categories/{id} [put]
func (c *categoryController) UpdateCategory(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idUint := conv.StringToUint(id)

	var req request.CreateCategoryRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Errorf("[CategoryController] UpdateCategory - 1: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[CategoryController] UpdateCategory - 2: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reqModel := model.Category{
		ID: idUint,
		Name: req.Name,
		Tagline: req.Tagline,
		Photo: req.Photo,
	}

	if err := c.categoryUsecase.UpdateCategory(ctx.Context(), &reqModel); err != nil {
		log.Errorf("[CategoryController] UpdateCategory - 3: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update category",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category updated successfully",
	})
}

func NewCategoryController(categoryUsecase usecase.CategoryUsecaseInterface) CategoryControllerInterface {
	return &categoryController{
		categoryUsecase: categoryUsecase,
	}
}
