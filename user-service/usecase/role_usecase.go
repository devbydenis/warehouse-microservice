package usecase

import (
	"context"
	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/repository"
)

type RoleUsecaseInterface interface{
	CreateRole(ctx context.Context, role model.Role) error
	UpdateRole(ctx context.Context, role model.Role) error
	DeleteRole(ctx context.Context, id uint) error
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	GetAllRoles(ctx context.Context) ([]model.Role, error)
}

type roleUsecase struct {
	roleRepo repository.RoleRepositoryInterface
}

// CreateRole implements RoleUsecaseInterface.
func (r *roleUsecase) CreateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.CreateRole(ctx, role)
}

// UpdateRole implements RoleUsecaseInterface.
func (r *roleUsecase) UpdateRole(ctx context.Context, role model.Role) error {
	return r.roleRepo.UpdateRole(ctx, role)
}

// DeleteRole implements RoleUsecaseInterface.
func (r *roleUsecase) DeleteRole(ctx context.Context, id uint) error {
	return r.roleRepo.DeleteRole(ctx, id)
}

// GetRoleByID implements RoleUsecaseInterface.
func (r *roleUsecase) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, id)
}

// GetAllRoles implements RoleUsecaseInterface.
func (r *roleUsecase) GetAllRoles(ctx context.Context) ([]model.Role, error) {
	return r.roleRepo.GetAllRoles(ctx)
}

func NewRoleUsecase(roleRepo repository.RoleRepositoryInterface) RoleUsecaseInterface {
	return &roleUsecase{roleRepo: roleRepo}
}