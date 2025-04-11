package models

type UpdateUserRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type CreateOrderRequest struct {
	Items []OrderItem `json:"items"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending processing shipped delivered cancelled"`
}

type DashboardStats struct {
	TotalOrders   int64 `json:"total_orders"`
	TotalProducts int64 `json:"total_products"`
	TotalUsers    int64 `json:"total_users"`
}

type CreateProductRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	Stock       int      `json:"stock" validate:"required,gte=0"`
	Category    string   `json:"category"`
	Images      []string `json:"images"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
	Category    string  `json:"category,omitempty"`
}

type CreateVendorRequest struct {
	User User `json:"user"`
}

type CreateStoreRequest struct {
	StoreName            string `json:"store_name" validate:"required"`
	StoreDescription     string `json:"store_description" validate:"required"`
	StoreLogo            string `json:"store_logo" validate:"required"`
	StoreUrl             string `json:"store_url" validate:"required"`
	StoreAddress         string `json:"store_address" validate:"required"`
	StoreWhatsappContact string `json:"store_whatsapp_contact" validate:"required"`
}
