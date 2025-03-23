package models

import "time"

type Vendor struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex" json:"user_id"`
	StoreName        string    `json:"store_name"`
	StoreDescription string    `json:"store_description"`
	ImageURL         string    `json:"image_url"`
	IsActive         bool      `json:"is_active"`
	Products         []Product `gorm:"foreignKey:VendorID" json:"products,omitempty"`
	Services         []Service `gorm:"foreignKey:VendorID" json:"services,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	User             User      `gorm:"foreignKey:UserID" json:"user"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VendorID    uint      `json:"vendor_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Service struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VendorID    uint      `json:"vendor_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Rate        float64   `json:"rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
