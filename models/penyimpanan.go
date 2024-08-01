package models

import "time"

// Penyimpanan model
type Penyimpanan struct {
	ID                 int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Jenis              string          `gorm:"type:varchar(255)" json:"jenis"`
	BiayaPenyimpanan   Float64OrString `gorm:"type:decimal(20,2)" json:"biaya_penyimpanan"`
	TanggalPenyimpanan time.Time       `json:"tanggal_penyimpanan"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Penyimpanan) TableName() string {
	return "penyimpanan"
}
