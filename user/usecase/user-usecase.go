package usecase

import (
	"errors"
	"fmt"
	"time"
	"user-service/domain"

	"github.com/gin-gonic/gin"
)

type userUsecase struct {
	userRepository      domain.UserRepository
	userCacheRepository domain.UserCacheRepository
}

func NewUserUsecase(ur domain.UserRepository, ucr domain.UserCacheRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository:      ur,
		userCacheRepository: ucr,
	}
}

func (u *userUsecase) Signup(ctx *gin.Context, ur *domain.User) (domain.User, error) {
	user, err := u.userRepository.Create(ctx, ur)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *userUsecase) Login(ctx *gin.Context, Username string, Password string) (domain.User, error) {
	count, err := u.userCacheRepository.GetByUsernameLogin(ctx, Username)
	if count >= 5 {
		return domain.User{}, errors.New("login failed, please try again later")
	}

	user, err := u.userRepository.GetByUsername(ctx, Username)
	if err != nil {
		return domain.User{}, fmt.Errorf("error retrieving user by username: %v", err)
	}

	if errPassword := user.CheckPassword(Password); errPassword != nil {

		if _, err := u.userCacheRepository.Incr(ctx, Username); err != nil {

			fmt.Println("Error incrementing failed login count:", err)
		}
		return domain.User{}, errPassword
	}
	if err := u.userCacheRepository.Delete(ctx, Username); err != nil {
		fmt.Println("Error deleting cache for user:", err)

	}

	return user, nil
}

// Update implements domain.UserUsecase.
func (u *userUsecase) Update(ctx *gin.Context, ur *domain.User) (domain.User, error) {
	user, err := u.userRepository.Update(ctx, ur)
	if err != nil {
		return user, err
	}
	return user, nil
}
