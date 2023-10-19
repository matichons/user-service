package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID         int        `db:"UserID" json:"userID"`
	Username       string     `db:"Username" json:"username"`
	Password       string     `db:"Password" json:"password"`
	Email          string     `db:"Email" json:"email"`
	ProfilePicture *string    `db:"ProfilePicture" json:"profilePicture,omitempty"`
	DateJoined     time.Time  `db:"DateJoined" json:"dateJoined"`
	LastLogin      *time.Time `db:"LastLogin" json:"lastLogin,omitempty"`
}

type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	GetByTitle(ctx context.Context, title string) (User, error)
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
	Create(ctx context.Context, ar *User) error
}

type UserRepository interface {
	Create(ctx context.Context, ar *User) error
	Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	GetByID(ctx context.Context, id int) (User, error)
	GetByUsername(ctx context.Context, Username string) (User, error)
	Update(ctx context.Context, ar *User) error
	Store(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
