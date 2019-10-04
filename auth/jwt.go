package auth

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"sinistra/go-gin-api/models"
)

func GenerateToken(user models.User) (string, error) {
	var err error
	secret := os.Getenv("JWT_SECRET")
	issuer := os.Getenv("JWT_ISSUER")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"iss":      issuer,
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
