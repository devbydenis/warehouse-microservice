package database

import (
	"micro-warehouse/user-service/model"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func SeedRole(db *gorm.DB) {
	roles := []model.Role{
		{Name: "Manager"},
		{Name: "Keeper"},
	}
	
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			log.Errorf("[RoleSeeder] SeedRole - 1: %v", err)
		} else {
			log.Infof("[RoleSeeder] SeedRole - 2: %v", "Role created successfully")
		}
	}
}