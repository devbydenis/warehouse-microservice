package controller

import (
	request "micro-warehouse/user-service/controller/request"
	response "micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"
	"micro-warehouse/user-service/pkg/pagination"
	"micro-warehouse/user-service/pkg/validator"
	"micro-warehouse/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserControllerInterface interface {
	CreateUser(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error

	GetUserByRoleName(c *fiber.Ctx) error

	AssignUserToRole(c *fiber.Ctx) error
	EditAssignUserToRole(c *fiber.Ctx) error
	GetUserRoleByID(c *fiber.Ctx) error
	GetAllUserRoles(c *fiber.Ctx) error
}

type userController struct {
	userUsecase usecase.UserUsecaseInterface
}



// CreateUser implements UserControllerInterface.
// @Summary Create User
// @Description Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.CreateUserRequest true "Create User Request Body"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users [post]
// @Security Bearer
func (u *userController) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.CreateUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] CreateUser - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] CreateUser - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userModel := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	if err := u.userUsecase.CreateUser(ctx, userModel); err != nil {
		log.Errorf("[UserController] CreateUser - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

// DeleteUser implements UserControllerInterface.
// @Summary Delete User
// @Description Delete a user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{id} [delete]
// @Security Bearer
func (u *userController) DeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	userId := conv.StringToUint(id)

	if err := u.userUsecase.DeleteUser(ctx, userId); err != nil {
		log.Errorf("[UserController] DeleteUser - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// GetAllUsers implements UserControllerInterface.
// @Summary Get All Users
// @Description Get all users
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search keyword"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order (asc/desc)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users [get]
// @Security Bearer
func (u *userController) GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetAllUsersRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[UserController] GetAllUser - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] GetAllUser - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := u.userUsecase.GetAllUsers(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[UserController] GetAllUser - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserResponse{}
	for _, user := range users {
		roleName := ""
		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		resp = append(resp, response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Photo:    user.Photo,
			RoleName: roleName,
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(total))

	response := response.GetAllUsersResponse{
		Users:      resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "Users fetched successfully",
	})

}

// GetUserByID implements UserControllerInterface.
// @Summary Get User By ID
// @Description Get a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{id} [get]
// @Security Bearer
func (u *userController) GetUserByID(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	userID := conv.StringToUint(id)

	user, err := u.userUsecase.GetUserByID(ctx, userID)
	if err != nil {
		log.Errorf("[UserController] GetUserByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	roleName := ""
	if len(user.Roles) > 0 {
		roleName = user.Roles[0].Name
	}

	response := response.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Photo:    user.Photo,
		RoleName: roleName,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

// UpdateUser implements UserControllerInterface.
// @Summary Update User
// @Description Update a user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/{id} [put]
// @Security Bearer
func (u *userController) UpdateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	var req request.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] UpdateUser - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] UpdateUser - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userID := conv.StringToUint(id)
	userModel := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Photo:    req.Photo,
		ID:       userID,
	}

	if req.Password != "" {
		hashedPassword, err := conv.HashPassword(req.Password)
		if err != nil {
			log.Errorf("[UserController] UpdateUser - 3: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		userModel.Password = hashedPassword
	}

	if err := u.userUsecase.UpdateUser(ctx, userModel); err != nil {
		log.Errorf("[UserController] UpdateUser - 4: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

// AssignUserToRole implements UserControllerInterface.
// @Summary Assign User To Role
// @Description Assign a user to a role
// @Tags Assign Role
// @Accept json
// @Produce json
// @Param request body request.AssignUserToRoleRequest true "Assign User To Role Request Body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/assign-role [post]
// @Security Bearer
func (u *userController) AssignUserToRole(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.AssignUserToRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] AssignUserToRole - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] AssignUserToRole - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := u.userUsecase.AssignUserToRole(ctx, req.UserID, req.RoleID); err != nil {
		log.Errorf("[UserController] AssignUserToRole - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User assigned to role successfully",
	})
}

// GetUserByRoleName implements UserControllerInterface.
// @Summary Get User By Role Name
// @Description Get a user by role name
// @Tags Assign Role
// @Accept json
// @Produce json
// @Param roleName path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/users/role/{roleName} [get]
// @Security Bearer
func (u *userController) GetUserByRoleName(c *fiber.Ctx) error {
	ctx := c.Context()
	roleName := c.Params("roleName")

	users, err := u.userUsecase.GetUserByRoleName(ctx, roleName)
	if err != nil {
		log.Errorf("[UserController] GetUserByRoleName - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserResponse{}
	for _, user := range users {
		roleName := ""
		if len(user.Roles) > 0 {
			roleName = user.Roles[0].Name
		}

		resp = append(resp, response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Photo:    user.Photo,
			RoleName: roleName,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Users fetched successfully",
	})
}

// GetUserRoleByID implements UserControllerInterface.
// @Summary Get User Role By ID
// @Description Get a user role by ID
// @Tags Assign Role
// @Accept json
// @Produce json
// @Param userRoleID path string true "User Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/assign-role/{userRoleID} [get]
// @Security Bearer
func (u *userController) GetUserRoleByID(c *fiber.Ctx) error {
	ctx := c.Context()
	userRoleIDStr := c.Params("userRoleID")
	userRoleID := conv.StringToUint(userRoleIDStr)

	userRole, err := u.userUsecase.GetUserRoleByID(ctx, userRoleID)
	if err != nil {
		log.Errorf("[UserController] GetUserRoleByID - 1: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    userRole,
		"message": "User role retrieved successfully",
	})
}

// EditAssignUserToRole implements UserControllerInterface.
// @Summary Edit Assign User To Role
// @Description Edit a user to a role
// @Tags Assign Role
// @Accept json
// @Produce json
// @Param request body request.AssignUserToRoleRequest true "Assign User To Role Request Body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/assign-role/{userRoleID} [put]
// @Security Bearer
func (u *userController) EditAssignUserToRole(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.AssignUserToRoleRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[UserController] EditAssignUserToRole - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] EditAssignUserToRole - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userRoleIDStr := c.Params("userRoleID")
	userRoleID := conv.StringToUint(userRoleIDStr)

	if err := u.userUsecase.EditAssignUserToRole(ctx, userRoleID, req.UserID, req.RoleID); err != nil {
		log.Errorf("[UserController] EditAssignUserToRole - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User role updated successfully",
	})
}

// GetAllUserRoles implements UserControllerInterface.
// @Summary Get All User Roles
// @Description Get all users with their roles
// @Tags Assign Role
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param search query string false "Search keyword"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order (asc/desc)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/assign-role [get]
// @Security Bearer
func (u *userController) GetAllUserRoles(c *fiber.Ctx) error {
	ctx := c.Context()

	var req request.GetAllUsersRequest
	if err := c.QueryParser(&req); err != nil {
		log.Errorf("[UserController] GetAllUserRoles - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[UserController] GetAllUserRoles - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := u.userUsecase.GetAllUserRoles(ctx, req.Page, req.Limit, req.Search, req.SortBy, req.SortOrder)
	if err != nil {
		log.Errorf("[UserController] GetAllUserRoles - 3: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	resp := []response.UserRoleResponse{}
	for _, user := range users {
		resp = append(resp, response.UserRoleResponse{
			ID:     user.ID,
			UserID: user.UserID,
			RoleID: user.RoleID,
			User: response.UserResponse{
				ID: user.User.ID,
				Name: user.User.Name,
				Email: user.User.Email,
				Phone: user.User.Phone,
				Photo: user.User.Photo,
			},
			Role: response.RoleResponse{
				ID:   user.Role.ID,
				Name: user.Role.Name,
			},
		})
	}

	paginationInfo := pagination.CalculatePagination(req.Page, req.Limit, int(total))

	response := response.GetAllUserRolesResponse{
		UserRoles:  resp,
		Pagination: paginationInfo,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    response,
		"message": "User roles fetched successfully",
	})

}

func NewUserController(userUsecase usecase.UserUsecaseInterface) UserControllerInterface {
	return &userController{
		userUsecase: userUsecase,
	}
}
