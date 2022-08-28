package models

import "gorm.io/gorm"

type Audience struct {
	gorm.Model
	AuthID   uint
	Auth     *Auth
	DomainID uint
	Domain   *Domain
}
