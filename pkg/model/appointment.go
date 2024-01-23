package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	DoctorID  uuid.UUID
	PatientID uuid.UUID
	CreatedAt time.Time `gorm:"autoCreateTime"`
	TimeStart time.Time
	TimeEnd   time.Time
}
