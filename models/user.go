package models

import "time"

// User model
type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Nama      string    `gorm:"type:varchar(255)" json:"nama"`
	Username  string    `gorm:"type:varchar(255);unique" json:"username"`
	Password  string    `gorm:"type:varchar(255)" json:"password"`
	Posisi    string    `gorm:"type:varchar(255)" json:"posisi"`
	HP        string    `gorm:"type:varchar(255)" json:"hp"`
	Alamat    string    `gorm:"type:text" json:"alamat"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (User) TableName() string {
	return "user"
}
