package auth

import (
	"github.com/gin-gonic/gin"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"golang.org/x/crypto/bcrypt"
)

// SignupBody describes the signup body with desired credentials
type SignupBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) Signup(ctx *gin.Context) (userID int, err error) {
	var body SignupBody
	err = ctx.ShouldBind(&body)
	if err != nil {
		return
	}

	hashedPassword, err := hashAndSaltPassword(body.Password)
	if err != nil {
		return
	}

	var user *m.User
	user, err = CreateUser(ctx, body.Name, body.Email, hashedPassword)
	if err != nil {
		return
	}

	return int(user.ID), nil
}

func hashAndSaltPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
