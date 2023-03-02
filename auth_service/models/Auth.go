package models

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	Email    string   `gorm:"unique;notNull;index"`
	Password string   `gorm:"notNull"`
	Salt     string   `gorm:"notNull;default:0"`
	Flags    []Flag   `gorm:"many2many:auths_flags;"`
	Domains  []Domain `gorm:"many2many:auths_domains;"`
}
