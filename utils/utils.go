package utils

import (
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
