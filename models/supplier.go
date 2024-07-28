package models

import "time"

// Supplier model
type Supplier struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama       string    `gorm:"type:varchar(255)" json:"nama"`
	Perusahaan string    `gorm:"type:varchar(255)" json:"perusahaan"`
	Kontak     string    `gorm:"type:varchar(255)" json:"kontak"`
	Alamat     string    `gorm:"type:text" json:"alamat"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Supplier) TableName() string {
	return "supplier"
}
