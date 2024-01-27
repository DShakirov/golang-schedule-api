package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type MedicalRecord struct {
	gorm.Model
	DoctorID     uuid.UUID
	DoctorEmail  string
	PatientID    uuid.UUID
	PatientEmail string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Text         string
}
