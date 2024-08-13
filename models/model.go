package models

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

type Float64OrString float64
type IntOrString int

func (i *IntOrString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*i = IntOrString(intValue)
		return nil
	}

	// If unmarshalling as int fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		intValue, err := strconv.Atoi(stringValue)
		if err != nil {
			return err
		}
		*i = IntOrString(intValue)
		return nil
	}

	return errors.New("cannot unmarshal to IntOrString")
}

func (f *Float64OrString) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64
	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err == nil {
		*f = Float64OrString(floatValue)
		return nil
	}

	// If unmarshalling as float64 fails, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		floatValue, err := strconv.ParseFloat(stringValue, 64)
		if err != nil {
			return err
		}
		*f = Float64OrString(floatValue)
		return nil
	}

	return errors.New("cannot unmarshal to Float64OrString")
}

type PemesananWithBarangWithSupplier struct {
	ID                  int             `json:"id"`
	IDUser              int             `gorm:"index" json:"id_user"`
	IDBarang            int             `gorm:"index" json:"id_barang"`
	BarangNama          string          `json:"barang_nama"`   // Tambahan
	SupplierNama        string          `json:"supplier_nama"` // Tambahan
	Kuantitas           int             `json:"kuantitas"`
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

type PenjualanWithBarangWithSupplier struct {
	ID               int             `json:"id"`
	IDUser           int             `gorm:"index" json:"id_user"`
	IDBarang         int             `gorm:"index" json:"id_barang"`
	BarangNama       string          `json:"barang_nama"`   // Tambahan
	SupplierNama     string          `json:"supplier_nama"` // Tambahan
	Kuantitas        int             `json:"kuantitas"`
	HargaSatuan      Float64OrString `gorm:"type:decimal(20,2)" json:"harga_satuan"`
	TotalHarga       Float64OrString `gorm:"type:decimal(20,2)" json:"total_harga"`
	TanggalPenjualan time.Time       `json:"tanggal_penjualan"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type EoqWithBarang struct {
	ID                 int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IDBarang           int       `gorm:"index" json:"id_barang"` // Index to reference barang
	NamaBarang         string    `json:"nama_barang"`            //tambahan
	NilaiEOQ           float64   `gorm:"type:decimal(10,2)" json:"nilai_eoq"`
	Periode            string    `gorm:"type:varchar(255)" json:"periode"`
	TanggalPerhitungan time.Time `json:"tanggal_perhitungan"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StockBarangModel struct {
	IDBarang                int    `json:"id_barang"`
	NamaBarang              string `gorm:"index" json:"nama_barang"`
	TotalKuantitasPemesanan int    `json:"total_kuantitas_pemesanan"`
	TotalKuantitasPenjualan int    `json:"total_kuantitas_penjualan"`
	StockBarang             int    `json:"stok_barang"`
}

type ReqStock struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
