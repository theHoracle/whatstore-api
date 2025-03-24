package models

import "time"

type Vendor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"user_id"`
	IsActive  bool      `json:"is_active"`
	Stores    []Store   `gorm:"foreignKey:VendorID" json:"stores,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Store struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VendorID    uint      `json:"vendor_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StoreLogo   string    `json:"store_logo"`
	Products    []Product `gorm:"foreignKey:StoreID" json:"products,omitempty"`
	Services    []Service `gorm:"foreignKey:StoreID" json:"services,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StoreID      uint      `json:"store_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency" gorm:"default:NGN"`
	SearchVector string    `gorm:"type:tsvector;index:idx_products_search,type:gin" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Service struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	StoreID      uint      `json:"store_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	Rate         float64   `json:"rate"`
	Currency     string    `json:"currency" gorm:"default:NGN"`
	SearchVector string    `gorm:"type:tsvector;index:idx_services_search,type:gin" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
