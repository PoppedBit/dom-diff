package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Id           uuid.UUID `gorm:"type:char(36);primaryKey"`
	Url          string    `gorm:"type:varchar(255);not null"`
	ItemSelector string    `gorm:"type:varchar(255);not null"`
	TextSelector string    `gorm:"type:varchar(255);not null"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) (err error) {
	j.Id = uuid.New()
	return
}

func (j *Job) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Exec("DELETE FROM runs WHERE job_id = ?", j.Id)
	return nil
}
