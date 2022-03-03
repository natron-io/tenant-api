package models

import "gorm.io/gorm"

type Tenant struct {
	gorm.Model
	Id             int32  `gorm:"primary_key"`
	GitHubTeamSlug string `gorm:"not null;unique"`
}
