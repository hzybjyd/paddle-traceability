package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:20;not null;index" json:"role"`
	PublicKey    string    `gorm:"type:text" json:"public_key,omitempty"`
	CompanyName  string    `gorm:"size:100" json:"company_name,omitempty"`
	Phone        string    `gorm:"size:20" json:"phone,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
