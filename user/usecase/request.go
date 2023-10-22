package usecase

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignupPayload struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required"`
	ProfilePicture string `json:"profilePicture"`
}

type UserUpdatePayload struct {
	Password       string `json:"password"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}
