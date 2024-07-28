package models

import "time"

// Penjualan model
type Penjualan struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IDUser           int       `gorm:"index" json:"id_user"`   // Index to reference user
	IDBarang         int       `gorm:"index" json:"id_barang"` // Index to reference barang
	Kuantitas        int       `json:"kuantitas"`
	HargaSatuan      float64   `gorm:"type:decimal(10,2)" json:"harga_satuan"`
	TotalHarga       float64   `gorm:"type:decimal(10,2)" json:"total_harga"`
	TanggalPenjualan time.Time `json:"tanggal_penjualan"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Penjualan) TableName() string {
	return "penjualan"
}
