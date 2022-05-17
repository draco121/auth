package utils

import (
	"auth/startup"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Utils struct {
}

func (u *Utils) GenerateHashedPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(pass), err
}

func (u *Utils) ComparePassword(hashedpassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	return err
}

func (u *Utils) CreateJwt(username string) (string, error) {
	var mySigningKey = []byte(startup.Config.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *Utils) ValidateJwt(token string) (string, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("an error occurred while parsing token")
		}
		return []byte(startup.Config.SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		return fmt.Sprint(claims["username"]), nil
	} else {
		return "", fmt.Errorf("session expired")
	}
}
