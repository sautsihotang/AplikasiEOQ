package models

import "time"

// User model
type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Nama      string    `gorm:"type:varchar(255)"`
	Username  string    `gorm:"type:varchar(255);unique"`
	Password  string    `gorm:"type:varchar(255)"`
	Posisi    string    `gorm:"type:varchar(255)"`
	HP        string    `gorm:"type:varchar(255)"`
	Alamat    string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
