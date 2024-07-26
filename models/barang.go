package models

import (
	"time"
)

// Barang model
type Barang struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	IDSupplier int       `gorm:"index"` // Index to reference supplier
	NamaBarang string    `gorm:"type:varchar(255)"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
