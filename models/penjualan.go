package models

import "time"

// Penjualan model
type Penjualan struct {
	ID               int `gorm:"primaryKey;autoIncrement"`
	IDUser           int `gorm:"index"` // Index to reference user
	IDBarang         int `gorm:"index"` // Index to reference barang
	Kuantitas        int
	HargaSatuan      float64 `gorm:"type:decimal(10,2)"`
	TotalHarga       float64 `gorm:"type:decimal(10,2)"`
	TanggalPenjualan time.Time
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
}
