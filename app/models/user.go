package models

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClerkUserID string    `gorm:"uniqueIndex" json:"clerk_user_id"`
	Name        string    `json:"name"`
	Email       string    `gorm:"uniqueIndex" json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Vendor      *Vendor   `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"vendor,omitempty"`
	Buyer       *Buyer    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"buyer,omitempty"`
}

type Vendor struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex" json:"user_id"`
	StoreName        string    `json:"store_name"`
	StoreDescription string    `json:"store_description"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Buyer struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex" json:"user_id"`
	PreferredPayment string    `json:"preferred_payment"`
	ShippingAddress  string    `json:"shipping_address"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
