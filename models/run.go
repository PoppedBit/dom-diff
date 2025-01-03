package models

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Run struct {
	gorm.Model
	Id      uuid.UUID `gorm:"type:char(36);primaryKey"`
	JobId   uuid.UUID `gorm:"type:char(36);not null"`
	Job     Job       `gorm:"foreignKey:JobId"`
	Matches int       `gorm:"type:integer;not null;default:0"`
}

func (r *Run) BeforeCreate(tx *gorm.DB) (err error) {
	r.Id = uuid.New()
	return
}

func (r *Run) BeforeDelete(tx *gorm.DB) (err error) {
	outputDir := filepath.Join(os.Getenv("OUTPUT_DIR"), r.JobId.String(), r.Id.String())
	os.RemoveAll(outputDir)
	return nil
}
