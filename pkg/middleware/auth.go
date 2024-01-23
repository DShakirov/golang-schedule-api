package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve token string from header
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Return secret key
			return []byte(os.Getenv("JWT_SECRET")), nil //
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"].(string) // Access the "user_id" field from claims
		// Use the userID variable here
		id := uuid.Must(uuid.FromString(userID))
		c.Set("uuid", id)
		isDoctor := claims["is_doctor"].(bool) // Acess the "is_doctor" field from claims
		c.Set("isDoctor", isDoctor)
	}
}
