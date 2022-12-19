package errorformat

import (
	"errors"
	"strings"
)

func ErrorMessage(err string) error {

	if strings.Contains(err, "pkey") {
		return errors.New("User ID already exist")
	} else if strings.Contains(err, "email_key") {
		return errors.New("Email already exist")
	} else if strings.Contains(err, "user not found") {
		return errors.New("Email is not registered")
	} else if strings.Contains(err, "hashedPassword") {
		return errors.New("Password is incorrect")
	}
	
	return errors.New(err)
}
