package c_model

import (
	"time"

	"gorm.io/gorm"
)

type GormModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Model struct {
	GormModel
}
