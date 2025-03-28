package models

import "time"

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusSuccess  OrderStatus = "success"
	OrderStatusRejected OrderStatus = "rejected"
)

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `json:"user_id"`
	StoreID     uint        `json:"store_id"`
	Status      OrderStatus `gorm:"type:string;default:'pending'" json:"status"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"items"`
	PaymentID   *string     `json:"payment_id,omitempty"` // For escrow/payment tracking
	User        User        `gorm:"foreignKey:UserID" json:"-"`
	Store       Store       `gorm:"foreignKey:StoreID" json:"-"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"` // Price at time of order
	Product   Product   `json:"product"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
