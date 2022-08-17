package models

import "gorm.io/gorm"

type Auth struct {
	gorm.Model
	UserId uint64
	Roles  []Role `gorm:"many2many:auths_roles;"`
	Flags  []Flag `gorm:"many2many:auths_flags;"`
}
