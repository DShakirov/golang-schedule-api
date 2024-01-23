package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	Type      string
	UserID    uuid.UUID
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Text      string
}
