package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zayn1510/goarchi/config"
)

const (
	ErrMissingToken  = "Missing token"
	ErrInvalidFormat = "Invalid token format"
	ErrInvalidToken  = "Invalid token"
	ErrInvalidClaims = "Invalid claims"
	ErrTokenExpired  = "Token has expired"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GetExpiredDuration() time.Duration {
	return config.JWTExpired
}
func getEnv(key, defaultValue string) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		val = defaultValue
	}
	// Konversi string ke int
	result, err := strconv.Atoi(val)
	if err != nil {
		// Jika gagal konversi, gunakan default (5)
		return 5
	}
	return result
}
func validateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, fmt.Errorf(ErrInvalidToken)
	}

	// Baru akses claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf(ErrInvalidClaims)
	}

	// Expired sudah otomatis dicek jwt.Parse, tapi boleh tetap eksplisit
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, nil, fmt.Errorf(ErrTokenExpired)
		}
	}

	return token, claims, nil
}

// Middleware JWT untuk Gin
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrMissingToken})
			c.Abort()
			return
		}

		// Token harus dalam format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidFormat})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// Validasi token
		_, claims, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Simpan username ke context untuk digunakan di handler
		c.Set("username", claims["username"])
		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}

func GenerateJWT(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(config.JWTExpired).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
