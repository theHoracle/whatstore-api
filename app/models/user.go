package models

import "time"

type User struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	ClerkID     string       `gorm:"uniqueIndex" json:"clerk_id"`
	Name        string       `json:"name"`
	Email       string       `gorm:"uniqueIndex" json:"email"`
	Username    string       `gorm:"uniqueIndex" json:"username"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Vendor      *Vendor      `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"vendor,omitempty"`
	UserDetails *UserDetails `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user_details,omitempty"`
	Orders      []Order      `json:"orders,omitempty"`
}

type UserDetails struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex" json:"user_id"`
	PreferredPayment string    `json:"preferred_payment"`
	ShippingAddress  string    `json:"shipping_address"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
