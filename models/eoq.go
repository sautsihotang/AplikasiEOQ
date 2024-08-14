package models

import "time"

// Eoq model
type Eoq struct {
	ID                 int             `gorm:"primaryKey;autoIncrement" json:"id"`
	IDBarang           int             `gorm:"index" json:"id_barang"` // Index to reference barang
	NilaiEOQ           float64         `gorm:"type:decimal(10,2)" json:"nilai_eoq"`
	D                  IntOrString     `json:"d"`
	S                  IntOrString     `json:"s"`
	H                  Float64OrString `gorm:"type:decimal(20,2)" json:"h"`
	Periode            string          `gorm:"type:varchar(255)" json:"periode"`
	TanggalPerhitungan time.Time       `json:"tanggal_perhitungan"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (Eoq) TableName() string {
	return "eoq"
}
