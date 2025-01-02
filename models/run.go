package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Run struct {
	gorm.Model
	Id      uuid.UUID `gorm:"type:char(36);primaryKey"`
	JobId   string    `gorm:"type:char(36);not null"`
	Job     Job       `gorm:"foreignKey:JobId"`
	Matches int       `gorm:"type:integer;not null;default:0"`
}

func (r *Run) BeforeCreate(tx *gorm.DB) (err error) {
	r.Id = uuid.New()
	return
}
