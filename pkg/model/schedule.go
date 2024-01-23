package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	DoctorID  uuid.UUID
	TimeStart time.Time
	TimeEnd   time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
