package auth

import (
	"time"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"

	"github.com/sumup-oss/go-pkgs/errors"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtConfig config.JWTConfig

func InitializeAuthConfig(cfg config.JWTConfig) {
	jwtConfig = cfg
}

func GenerateJWT(id, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expiration time: 24 hours
		"iat":   time.Now().Unix(),                     // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string) (string, string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, idOk := claims["id"].(string)
		email, emailOk := claims["email"].(string)

		if !idOk || !emailOk {
			return "", "", errors.New("invalid claims format")
		}

		return id, email, nil
	}

	return "", "", errors.New("invalid token")
}
