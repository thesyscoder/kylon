package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cluster struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	Name       string `gorm:"uniqueIndex;not null"`
	Kubeconfig string `gorm:"type:text;not null"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
