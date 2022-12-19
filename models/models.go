package models

import (
	"errors"
	"html"
	"strings"
	"task-vix-btpns/app"
	"task-vix-btpns/helpers/hash"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

type User struct {
	ID        string    `gorm:"primary_key; unique" json:"id"`
	Username  string    `gorm:"size:255;not null;" json:"username"`
	Email     string    `gorm:"size:255;not null; unique" json:"email"`
	Password  string    `gorm:"size:255;not null;" json:"password"`
	Photos    Photo     `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL;" json:"photos"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Photo struct {
	ID       int       `gorm:"primary_key;auto_increment" json:"id"`
	Title    string    `gorm:"size:255;not null" json:"title"`
	Caption  string    `gorm:"size:255;not null" json:"caption"`
	PhotoUrl string    `gorm:"size:255;not null;" json:"photo_url"`
	UserID   string    `gorm:"not null" json:"user_id"`
	Owner    app.Owner `gorm:"owner"`
}

// USER METHODS

//Inisialize user data
func (u *User) Init() {
	u.ID = uuid.New().String()                                    //Generate new uuid
	u.Username = html.EscapeString(strings.TrimSpace(u.Username)) //Escape string
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

// Change password to hashed password
func (u *User) HashPassword() error {
	hashedPassword, err := hash.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Check password
func (u *User) CheckPassword(providedPassword string) error {
	err := hash.CheckPasswordHash(u.Password, providedPassword)
	if err != nil {
		return err
	}
	return nil
}

//Validate user data
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) { //Convert to lowercase

		case "login": //Login case
			if u.Email == "" {
				return errors.New("Email is required")
			}
			if u.Password == "" {
				return errors.New("Password is required")
			}
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Email is invalid")
			}
			return nil

		case "register": //Register case
			if u.ID == "" {
				return errors.New("ID is required")
			} else if u.Email == "" {
				return errors.New("Email is required")
			} else if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Email is invalid")
			} else if u.Username == "" {
				return errors.New("Username is required")
			} else if u.Password == "" {
				return errors.New("Password is required")
			} else if len(u.Password) < 8 {
				return errors.New("Password must be at least 8 characters")
			}

			return nil

		case "update": //Update case
			if u.ID == "" {
				return errors.New("ID is required")
			} else if u.Email == "" {
				return errors.New("Email is required")
			} else if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("invalid email")
			} else if u.Username == "" {
				return errors.New("Username is required")
			} else if u.Password == "" {
				return errors.New("Password is required")
			} else if len(u.Password) < 8 {
				return errors.New("Password must be at least 8 characters")
			}

			return nil

		default:
			return nil
	}
}

//PHOTO METHODS
//Function to initialize Photo data
func (p *Photo) Init() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title)) //Escape string
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl = html.EscapeString(strings.TrimSpace(p.PhotoUrl))
}

//Function to validate Photo data
func (p *Photo) Validate(action string) error {
	switch strings.ToLower(action) { //Convert to lowercase
		case "upload": // Create/Upload case
			if p.Title == "" {
				return errors.New("Title is required")
			} else if p.Caption == "" {
				return errors.New("Caption is required")
			} else if p.UserID == "" {
				return errors.New("UserID is required")
			}
			return nil

		case "change": //Change case
			if p.Title == "" {
				return errors.New("Title is required")
			} else if p.Caption == "" {
				return errors.New("Caption is required")
			} else if p.PhotoUrl == "" {
				return errors.New("PhotoUrl is required")
			}
			return nil

		default:
			return nil
	}
}
