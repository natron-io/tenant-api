package models

import "gorm.io/gorm"

type CPUCost struct {
	gorm.Model
	Value    float64 `gorm:"not null"`
	TenantId int32   `gorm:"not null"`
}

type MemoryCost struct {
	gorm.Model
	Value    float64 `gorm:"not null"`
	TenantId int32   `gorm:"not null"`
}

type IngressCost struct {
	gorm.Model
	Value    float64 `gorm:"not null"`
	TenantId int32   `gorm:"not null"`
}

type StorageCost struct {
	gorm.Model
	Value        float64 `gorm:"not null"`
	TenantId     int32   `gorm:"not null"`
	StorageClass string  `gorm:"not null"`
}
