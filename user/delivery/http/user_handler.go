package http

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"user-service/domain"
	"user-service/middlewares"
	"user-service/user/usecase"

	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UUsecase domain.UserUsecase
}

func NewUserHandler(g *gin.Engine, us domain.UserUsecase) {
	handler := &UserHandler{
		UUsecase: us,
	}

	g.POST("/signup", handler.Signup)
	g.POST("/login", handler.Login)
	protected := g.Group("").Use(middlewares.Authz("verysecretkey", "AuthService"))
	{
		protected.PATCH("/update-profile", handler.Update)
	}
}

func (u *UserHandler) Signup(g *gin.Context) {
	var input usecase.UserSignupPayload
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		return
	}
	user, errDb := u.UUsecase.Signup(g, &domain.User{
		Username:       input.Username,
		Password:       input.Password,
		Email:          input.Email,
		ProfilePicture: &input.ProfilePicture,
	})
	if errDb != nil {
		g.JSON(http.StatusConflict, ResponseError{Message: errDb.Error()})
		return
	}
	tokenResponse, err := GenerateToken(user)
	if err != nil {
		log.Println(err)
		g.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		g.Abort()
		return
	}
	g.JSON(http.StatusCreated, tokenResponse)
}

func (u *UserHandler) Login(g *gin.Context) {
	var payload usecase.LoginPayload
	if err := g.ShouldBindJSON(&payload); err != nil {
		g.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		return
	}
	user, errDb := u.UUsecase.Login(g, payload.Username, payload.Password)
	if errDb != nil {
		g.JSON(http.StatusConflict, ResponseError{Message: errDb.Error()})
		return
	}
	now := time.Now()
	errChan := make(chan error)

	go func() {
		_, err := u.UUsecase.Update(g, &domain.User{Id: user.Id, LastLogin: &now})
		errChan <- err
	}()

	if err := <-errChan; err != nil {
		fmt.Printf("%v ", err)
	}
	tokenResponse, err := GenerateToken(user)
	if err != nil {
		log.Println(err)
		g.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		g.Abort()
		return
	}
	g.JSON(200, tokenResponse)
}

func (u *UserHandler) Update(g *gin.Context) {
	var input usecase.UserUpdatePayload
	if err := g.ShouldBindJSON(&input); err != nil {
		g.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
		return
	}
	user, errDb := u.UUsecase.Update(g, &domain.User{
		Password:       input.Password,
		Email:          input.Email,
		ProfilePicture: &input.ProfilePicture,
	})
	if errDb != nil {
		g.JSON(http.StatusConflict, ResponseError{Message: errDb.Error()})
		return
	}
	tokenResponse, err := GenerateToken(user)
	if err != nil {
		g.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusCreated, tokenResponse)
}

func GenerateToken(user domain.User) (usecase.LoginResponse, error) {
	jwtWrapper := middlewares.JwtWrapper{
		SecretKey:         "verysecretkey",
		Issuer:            "AuthService",
		ExpirationMinutes: 3600,
		ExpirationHours:   12,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.Id, user.Username)
	if err != nil {
		return usecase.LoginResponse{}, err
	}
	signedtoken, err := jwtWrapper.RefreshToken(user.Id, user.Username)
	if err != nil {

		return usecase.LoginResponse{}, err
	}
	tokenResponse := usecase.LoginResponse{
		Token:        signedToken,
		RefreshToken: signedtoken,
	}
	return tokenResponse, nil
}
