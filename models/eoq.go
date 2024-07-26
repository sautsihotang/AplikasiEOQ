package models

import "time"

type Eoq struct {
	ID                 int     `gorm:"primaryKey;autoIncrement"`
	IDBarang           int     `gorm:"index"` // Index to reference barang
	NilaiEOQ           float64 `gorm:"type:decimal(10,2)"`
	Periode            string  `gorm:"type:varchar(255)"`
	TanggalPerhitungan time.Time
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}
