package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lits-06/manage-user/internal/models"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Set("isNewSession", true)
		c.Next()
		return
	}

	tokenString := authHeader[len("Bearer "):]
	var userClaim models.Claim
	token, err := jwt.ParseWithClaims(tokenString, &userClaim, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRETKEY")), nil
	})

	if err != nil || !token.Valid {
		c.Set("isNewSession", true)
		c.Next()
		return
	}

	c.Set("userEmail", userClaim.Email)
	c.Set("isNewSession", false)
	c.Next()
}