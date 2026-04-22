package model

import (
	"errors"
	"net/mail"
	"strings"
)

type Profile struct {
	UserID   int64  `db:"user_id" json:"user_id"`
	Name     string `db:"name" json:"name"`
	FullName string `db:"full_name" json:"full_name"`
	Email    string `db:"email" json:"email"`
	Bio      string `db:"bio" json:"bio"`
}

func (p *Profile) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	p.FullName = strings.TrimSpace(p.FullName)
	p.Email = strings.TrimSpace(p.Email)
	p.Bio = strings.TrimSpace(p.Bio)

	if p.UserID == 0 {
		return errors.New("user_id is required")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if len(p.Name) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if p.Email == "" {
		return errors.New("email is required")
	}

	if _, err := mail.ParseAddress(p.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}
