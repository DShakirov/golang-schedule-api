package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Prescription struct {
	gorm.Model
	DrugName     string
	Dosage       string
	Duration     time.Duration
	DoctorID     uuid.UUID
	DoctorEmail  string
	PatientID    uuid.UUID
	PatientEmail string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
