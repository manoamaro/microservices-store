package models

import (
	"github.com/manoamaro/microservices-store/commons/pkg/collections"
	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Email    string   `gorm:"unique;notNull;index"`
	Password string   `gorm:"notNull"`
	Salt     string   `gorm:"notNull;default:0"`
	Flags    []Flag   `gorm:"many2many:auths_flags;"`
	Domains  []Domain `gorm:"many2many:auths_domains;"`
}

func (receiver Auth) FlagsArray() []string {
	return collections.MapTo(receiver.Flags, func(i Flag) string {
		return i.Name
	})
}

func (receiver Auth) DomainArray() []string {
	return collections.MapTo(receiver.Domains, func(i Domain) string {
		return i.Domain
	})
}
