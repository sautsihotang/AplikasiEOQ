package models

import "time"

// Penyimpanan model
type Penyimpanan struct {
	ID                 int     `gorm:"primaryKey;autoIncrement"`
	Jenis              string  `gorm:"type:varchar(255)"`
	BiayaPenyimpanan   float64 `gorm:"type:decimal(10,2)"`
	TanggalPenyimpanan time.Time
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}
