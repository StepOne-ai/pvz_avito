package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("some_super_secret_key")

type jwtCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				c.JSON(http.StatusBadRequest, models.Error{Message: "Token missing in both cookies and Authorization header"})
				c.Abort()
				return
			}
			tokenString = authHeader[7:]
		}

		token, err := verifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Error verifying token"})
			c.Abort()
			return
		}

		fmt.Printf("Token verified successfully. Claims: %+v\n", token.Claims)

		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := c.Cookie("role")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "No cookies named role"})
			c.Abort()
			return
		}
		isValid := false

		for _, a := range allowedRoles {
			if role == a {
				isValid = true
				break
			}
		}

		if !isValid {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Role is not enough"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func JwtSecret() []byte {
	return jwtSecret
}

func GenerateToken(user models.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
