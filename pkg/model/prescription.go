package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Prescription struct {
	gorm.Model
	DrugName  string
	Dosage    string
	Duration  time.Duration
	DoctorID  uuid.UUID
	PatientID uuid.UUID
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
