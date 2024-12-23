package models

import "gorm.io/gorm"

type Run struct {
	gorm.Model
	Id    string `gorm:"type:char(36);primaryKey"`
	JobId string `gorm:"type:char(36);not null"`
	Job   Job    `gorm:"foreignKey:JobId"`
}
