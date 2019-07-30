package datamodels

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Users is our Users example model.
// Keep note that the tags for public-use (for our web app)
// should be kept in other file like "web/viewmodels/user.go"
// which could wrap by embedding the datamodels.Users or
// define completely new fields instead but for the shake
// of the example, we will use this datamodel
// as the only one Users model in our application.
type Users struct {
	ID              int64  `json:"id" form:"id"`
	Name            string `json:"name" form:"name"`
	Email           string
	EmailVerifiedAt time.Time
	Password        string
	RememberToken   string
	CreatedAt       time.Time `json:"created_at" form:"created_at"`
	UpdateAt        time.Time
}

// IsValid can do some very very simple "low-level" data validations.
func (u Users) IsValid() bool {
	return u.ID > 0
}

// GeneratePassword will generate a hashed password for us based on the
// user's input.
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
