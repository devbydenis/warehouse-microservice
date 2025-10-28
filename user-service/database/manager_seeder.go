package database

import (
	"log"

	"micro-warehouse/user-service/model"
	"micro-warehouse/user-service/pkg/conv"

	"gorm.io/gorm"
)

func SeedManager(db *gorm.DB) {
	bytes, err := conv.HashPassword("manager123")
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	modelRole := model.Role{}
	err = db.Where("name = ?", "Manager").First(&modelRole).Error
	if err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}

	admin := model.User{
		Name: "manager",
		Email: "manager@mail.com",
		Password: bytes,
		Roles: []model.Role{modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "manager@mail.com"}).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
	}	else {
		log.Printf("Admin %s created", admin.Name)
	}

}