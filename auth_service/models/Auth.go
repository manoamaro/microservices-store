package models

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	Email    string
	Password string
	Salt     string
	Flags    []Flag `gorm:"many2many:auths_flags;"`
}
