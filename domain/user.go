package domain

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id             int        `db:"id" json:"id" gorm:"primaryKey"`
	Username       string     `db:"Username" json:"username" gorm:"unique;not null"`
	Password       string     `db:"Password" json:"password"`
	Email          string     `db:"Email" json:"email" gorm:"unique;not null"`
	ProfilePicture *string    `db:"ProfilePicture" json:"profilePicture,omitempty"`
	DateJoined     time.Time  `db:"DateJoined" json:"dateJoined"`
	LastLogin      *time.Time `db:"LastLogin" json:"lastLogin,omitempty"`
}

type UserUsecase interface {
	// Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	// GetByID(ctx context.Context, id int64) (User, error)
	// Update(ctx context.Context, ar *User) error
	// GetByTitle(ctx context.Context, title string) (User, error)
	// Store(context.Context, *User) error
	// Delete(ctx context.Context, id int64) error
	// GetByUsername(ctx *gin.Context, Username string) (User, error)
	Signup(ctx *gin.Context, ur *User) (User, error)
	Login(ctx *gin.Context, Username string, Password string) (User, error)
	Update(ctx *gin.Context, ur *User) (User, error)
}

type UserRepository interface {
	Create(ctx *gin.Context, ur *User) (User, error)
	// Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	// GetByID(ctx context.Context, id int) (User, error)
	GetByUsername(ctx *gin.Context, Username string) (User, error)
	Update(ctx *gin.Context, ur *User) (User, error)
	// Store(ctx context.Context, a *User) error
	// Delete(ctx context.Context, id int64) error
}

type UserCacheRepository interface {
	Incr(ctx *gin.Context, key string) (int64, error)
	GetByUsernameLogin(ctx *gin.Context, key string) (int64, error)
	Delete(ctx *gin.Context, key string) error
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
