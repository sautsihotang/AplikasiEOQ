package models

import "time"

// Supplier model
type Supplier struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	Nama       string    `gorm:"type:varchar(255)"`
	Perusahaan string    `gorm:"type:varchar(255)"`
	Kontak     string    `gorm:"type:varchar(255)"`
	Alamat     string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
