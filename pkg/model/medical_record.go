package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type MedicalRecord struct {
	gorm.Model
	DoctorID  uuid.UUID
	PatientID uuid.UUID
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Text      string
}
