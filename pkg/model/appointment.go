package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	DoctorID     uuid.UUID
	DoctorEmail  string
	PatientID    uuid.UUID
	PatientEmail string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	TimeStart    time.Time
	TimeEnd      time.Time
}
