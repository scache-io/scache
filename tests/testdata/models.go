package models

// User 用户结构体 - 用于测试
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active"`
}

// Product 产品结构体 - 用于测试
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

// Order 订单结构体 - 用于测试复杂类型
type Order struct {
	ID         string  `json:"id"`
	UserID     int     `json:"user_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}
