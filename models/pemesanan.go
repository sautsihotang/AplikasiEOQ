package models

import "time"

// Pemesanan model
type Pemesanan struct {
	ID                  int `gorm:"primaryKey;autoIncrement"`
	IDUser              int `gorm:"index"` // Index to reference user
	IDBarang            int `gorm:"index"` // Index to reference barang
	Kuantitas           int
	HargaSatuan         float64 `gorm:"type:decimal(10,2)"`
	TotalHarga          float64 `gorm:"type:decimal(10,2)"`
	BiayaTelepon        float64 `gorm:"type:decimal(10,2)"`
	BiayaAdm            float64 `gorm:"type:decimal(10,2)"`
	BiayaTransportasi   float64 `gorm:"type:decimal(10,2)"`
	TotalBiayaPemesanan float64 `gorm:"type:decimal(10,2)"`
	TanggalPemesanan    time.Time
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}
