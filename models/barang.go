package models

import (
	"time"
)

// Barang model
type Barang struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IDSupplier int       `gorm:"index" json:"id_supplier"` // Index to reference supplier
	NamaBarang string    `gorm:"type:varchar(255)" json:"nama_barang"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Barang) TableName() string {
	return "barang"
}
