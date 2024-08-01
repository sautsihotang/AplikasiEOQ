package models

import "time"

// Pemesanan model
type Pemesanan struct {
	ID                  int             `gorm:"primaryKey;autoIncrement" json:"id"`
	IDUser              int             `gorm:"index" json:"id_user"`   // Index to reference user
	IDBarang            int             `gorm:"index" json:"id_barang"` // Index to reference barang
	Kuantitas           IntOrString     `json:"kuantitas"`
	HargaSatuan         Float64OrString `gorm:"type:decimal(20,2)" json:"harga_satuan"`
	TotalHarga          Float64OrString `gorm:"type:decimal(20,2)" json:"total_harga"`
	BiayaTelepon        Float64OrString `gorm:"type:decimal(20,2)" json:"biaya_telepon"`
	BiayaAdm            Float64OrString `gorm:"type:decimal(20,2)" json:"biaya_adm"`
	BiayaTransportasi   Float64OrString `gorm:"type:decimal(20,2)" json:"biaya_transportasi"`
	TotalBiayaPemesanan Float64OrString `gorm:"type:decimal(20,2)" json:"total_biaya_pemesanan"`
	TanggalPemesanan    time.Time       `json:"tanggal_pemesanan"`
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Pemesanan) TableName() string {
	return "pemesanan"
}
