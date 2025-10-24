package repository

import (
	"context"
	"errors"
	"micro-warehouse/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// USER: GET by email, GET by ID, Create, Update, Delete, GET by role name, GET all
// ROLE: GET by id, GET all, create, Update, Delete
// USER ROLE: Assign user to role, getUserRoleByID, get all user role, edit assign user to role

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetAllUsers(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id uint) error

	GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error)

	AssignUserToRole(ctx context.Context, userID uint, roleID uint) error
	EditAssignUserToRole(ctx context.Context, assignRoleID uint, userID uint, roleID uint) error
	GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error)
	GetAllUserRoles(ctx context.Context, page, limit int, search, sortBy, sortOrder string) ([]model.UserRole, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// AssignUserToRole implements UserRepositoryInterface.
func (u *userRepository) AssignUserToRole(ctx context.Context, userID uint, roleID uint) error {
	select {
		case <-ctx.Done():
			log.Errorf("[userRepository] AssignUserToRole - 1: %v", ctx.Err())
			return ctx.Err()
		default:
	}

	userRole := model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	return u.db.WithContext(ctx).Create(&userRole).Error

}

// CreateUser implements UserRepositoryInterface.
func (u *userRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	//
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] CreateUser - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
	}

	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Errorf("[UserRepository] CreateUser - 2: %v", err)
		return nil, err
	}

	if user.ID == 0 {
		log.Errorf("[UserRepository] CreateUser - 3: %v", "User ID is 0")
		return nil, errors.New("user ID is invalid after create")
	}

	return &user, nil

}

// DeleteUser implements UserRepositoryInterface.
func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] DeleteUser - 1: %v", ctx.Err())
			return ctx.Err()
		default:
	}

	modelUser := model.User{}

	if err := u.db.WithContext(ctx).
		Select("id", "name", "email", "password", "photo", "phone").
		Preload("Roles").
		Where("id = ?", id).
		First(&modelUser).Error; err != nil {
			log.Errorf("[UserRepository] DeleteUser - 2: %v", err)
			return err
	}

	return u.db.WithContext(ctx).Delete(&modelUser).Error
}

// EditAssignUserToRole implements UserRepositoryInterface.
func (u *userRepository) EditAssignUserToRole(ctx context.Context, assignRoleID uint, userID uint, roleID uint) error {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] EditAssignUserToRole - 1: %v", ctx.Err())
			return ctx.Err()
		default:
	}

	userRole := model.UserRole{}

	if err := u.db.WithContext(ctx).
		Select("id", "user_id", "role_id").
		Where("id = ?", assignRoleID).
		First(&userRole).Error; err != nil {
			log.Errorf("[UserRepository] EditAssignUserToRole - 2: %v", err)
			return err
	}

	userRole.UserID = userID
	userRole.RoleID = roleID

	return u.db.WithContext(ctx).Save(&userRole).Error
}

// GetAllUserRoles implements UserRepositoryInterface.
func (u *userRepository) GetAllUserRoles(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.UserRole, int64, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] GetAllUserRoles - 1: %v", ctx.Err())
			return nil, 0, ctx.Err()
		default:
	}

	userRoles := []model.UserRole{}
	var totalRecords int64

	// Build Query
	query := u.db.WithContext(ctx).Model(&model.UserRole{})

	// Apply search filter if provided
	if search != "" {
		query = query.Joins("JOIN users ON user_roles.user_id = users.id").
			Joins("JOIN roles ON user_roles.role_id = roles.id").
			Where("users.name ILIKE ? OR users.email ILIKE ? OR roles.name ILIKE ?", 
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total record count
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Errorf("[UserRepository] GetAllUserRoles - 2: %v", err)
		return nil, 0, err
	}

	// Apply sorting
	if sortBy == "" {
		if sortOrder == "" {
			sortOrder = "asc"
		}
		query = query.Order(sortBy + " " + sortOrder)
	} else {
		query = query.Order("id desc")
	}

	// apply pagination
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// Execute query with preloads
	if err := query.
		Preload("User").
		Preload("Role").
		Find(&userRoles).Error; err != nil {
			log.Errorf("[UserRepository] GetAllUserRoles - 3: %v", err)
			return nil, 0, err
	}

	return userRoles, totalRecords, nil
}

// GetAllUsers implements UserRepositoryInterface.
func (u *userRepository) GetAllUsers(ctx context.Context, page int, limit int, search string, sortBy string, sortOrder string) ([]model.User, int64, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] GetAllUsers - 1: %v", ctx.Err())
			return nil, 0, ctx.Err()
		default:
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// calculate offset
	offset := (page - 1) * limit

	// build query
	query := u.db.WithContext(ctx).Model(&model.User{})

	// Add search filter if provided
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		log.Errorf("[UserRepository] GetAllUsers - 2: %v", err)
		return nil, 0, err
	}

	// Get paginated data
	modelUsers := []model.User{}
	if err := query.Select("id", "name", "email", "password", "photo", "phone", "created_at").
		Preload("Roles").
		Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(limit).
		Find(&modelUsers).Error; err != nil {
			log.Errorf("[Repository] GetAllUsers - 2: %v", err)
			return nil, 0, err
	}

	return modelUsers, totalRecords, nil
}

// GetUserByEmail implements UserRepositoryInterface.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] GetUserByEmail - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
	}
	
	modelUsers := model.User{}
	if err := u.db.WithContext(ctx).Select("id", "name", "email", "password", "photo", "phone", "created_at").
		Where("email = ?", email).
		Preload("Roles").
		First(&modelUsers).Error; err != nil {
			log.Errorf("[UserRepository] GetUserByEmail - 2: %v", err)
			return nil, err
	}

	return &modelUsers, nil
}

// GetUserByID implements UserRepositoryInterface.
func (u *userRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
		select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] GetUserByID - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
	}
	
	modelUsers := model.User{}
	if err := u.db.WithContext(ctx).Select("id", "name", "email", "password", "photo", "phone", "created_at").
		Where("id = ?", id).
		Preload("Roles").
		First(&modelUsers).Error; err != nil {
			log.Errorf("[UserRepository] GetUserByID - 2: %v", err)
			return nil, err
	}

	return &modelUsers, nil
}

// GetUserByRoleName implements UserRepositoryInterface.
func (u *userRepository) GetUserByRoleName(ctx context.Context, roleName string) ([]model.User, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] GetUserByRoleName - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
	}

	users := []model.User{}

	// use subquery to get user ID's whos have the specified role
	subQuery := u.db.Table("user_role").
		Select("user_role.user_id").
		Joins("JOIN roles ON user_role.role_id = roles.id").
		Where("roles.name = ?", roleName)

	// main query with preload roles
	if err := u.db.WithContext(ctx).
		Where("id IN (?)", subQuery).
		Preload("Roles").
		Find(&users).Error; err != nil {
			log.Errorf("[UserRepository] GetUserByRoleName - 2: %v", err)
			return nil, err
	}

	return users, nil
}

// GetUserRoleByID implements UserRepositoryInterface.
func (u *userRepository) GetUserRoleByID(ctx context.Context, assignRoleID uint) (*model.UserRole, error) {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] GetUserRoleByID - 1: %v", ctx.Err())
			return nil, ctx.Err()
		default:
	}

	userRole := model.UserRole{}

	if err := u.db.WithContext(ctx).
		Select("id", "user_id", "role_id", "updated_at").
		Preload("User").
		Preload("Role").
		Where("id = ?", assignRoleID).
		First(&userRole).Error; err != nil {
			log.Errorf("[UserRepository] GetUserRoleByID - 2: %v", err)
			return nil, err
	}

	return &userRole, nil
}

// UpdateUser implements UserRepositoryInterface.
func (u *userRepository) UpdateUser(ctx context.Context, user model.User) error {
	select {
		case <-ctx.Done():
			log.Errorf("[UserRepository] UpdateUser - 1: %v", ctx.Err())
			return ctx.Err()
		default:
	}

	modelUser := model.User{}

	if err := u.db.WithContext(ctx).
		Select("id", "name", "email", "password", "photo", "phone").
		Preload("Roles").
		Where("id = ?", user.ID).
		First(&modelUser).Error; err != nil {
			log.Errorf("[UserRepository] UpdateUser - 2: %v", err)
			return err
	}

	modelUser.Name = user.Name
	modelUser.Email = user.Email
	if user.Password != "" {
		modelUser.Password = user.Password
	}
	modelUser.Photo = user.Photo
	modelUser.Phone = user.Phone

	return u.db.WithContext(ctx).Save(&modelUser).Error
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}
