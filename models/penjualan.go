package models

import "time"

// Penjualan model
type Penjualan struct {
	ID               int             `gorm:"primaryKey;autoIncrement" json:"id"`
	IDUser           int             `gorm:"index" json:"id_user"`   // Index to reference user
	IDBarang         int             `gorm:"index" json:"id_barang"` // Index to reference barang
	Kuantitas        IntOrString     `json:"kuantitas"`
	HargaSatuan      Float64OrString `gorm:"type:decimal(20,2)" json:"harga_satuan"`
	TotalHarga       Float64OrString `gorm:"type:decimal(20,2)" json:"total_harga"`
	TanggalPenjualan time.Time       `json:"tanggal_penjualan"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Penjualan) TableName() string {
	return "penjualan"
}
