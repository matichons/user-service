package repository

import (
	"user-service/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handlerUserRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) domain.UserRepository {
	return &handlerUserRepository{db}
}

func (u *handlerUserRepository) Create(ctx *gin.Context, user *domain.User) (domain.User, error) {
	err := user.HashPassword(user.Password)
	if err != nil {
		return *user, err
	}
	if result := u.DB.Create(&user); result.Error != nil {
		return *user, result.Error
	}
	return *user, nil
}

func (u *handlerUserRepository) GetByUsername(ctx *gin.Context, Username string) (domain.User, error) {
	var user domain.User
	if result := u.DB.First(&user, "username = ?", Username); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (u *handlerUserRepository) Update(ctx *gin.Context, ur *domain.User) (domain.User, error) {
	var user domain.User
	id := ctx.MustGet("id").(int)
	if err := u.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	if ur.Password != "" {
		ur.HashPassword(ur.Password)
	}

	if errUpdate := u.DB.Model(&user).Updates(&ur).Error; errUpdate != nil {
		return user, errUpdate
	}
	return user, nil
}
