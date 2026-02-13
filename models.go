package main

import "time"

// User represents a user account
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Role         string // "USER" or "ADMIN"
	CreatedAt    time.Time
}

// Product represents a merchandise product
type Product struct {
	ID               int
	Name             string
	Description      string
	Price            float64
	ImageURL         string
	HighlightText    string // Special highlight text (e.g., "Limited Drop", "Only Today")
	HighlightEnabled bool   // Whether to show the highlight
	CreatedAt        time.Time
}

// Cart represents a user's shopping cart
type Cart struct {
	ID        int
	UserID    int
	CreatedAt time.Time
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        int
	CartID    int
	ProductID int
	Size      string
	Quantity  int
	Product   *Product // Populated when needed
	AddedAt   time.Time
}

// Order represents a completed order
type Order struct {
	ID              int
	UserID          int
	UserName        string // For admin view
	TotalAmount     float64
	Status          string
	DeliveryAddress string
	PhoneNumber     string
	PaymentMethod   string
	CreatedAt       time.Time
	Items           []OrderItem // Populated when needed
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          int
	OrderID     int
	ProductID   int
	ProductName string
	Size        string
	Quantity    int
	Price       float64
	CreatedAt   time.Time
}

// Notification represents a user notification
type Notification struct {
	ID               int
	UserID           *int // NULL for global notifications
	NotificationType string
	Message          string
	IsGlobal         bool
	CreatedAt        time.Time
}

// Session represents user session data
type Session struct {
	UserID   int
	UserName string
	Email    string
	Role     string // "USER" or "ADMIN"
}

// SavedCard represents a user's saved payment card
// ⚠️ SIMULATED PAYMENT DATA - FOR EDUCATIONAL DEMO ONLY
type SavedCard struct {
	ID             int
	UserID         int
	CardholderName string
	CardNumber     string // ⚠️ Stored in plain text (demo only)
	ExpiryMonth    string
	ExpiryYear     string
	CVV            string // ⚠️ Stored in plain text (demo only)
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// CheckoutData represents data needed for checkout
type CheckoutData struct {
	DeliveryAddress string
	PhoneNumber     string
	CardholderName  string
	CardNumber      string
	ExpiryMonth     string
	ExpiryYear      string
	CVV             string
}

// AdminStats represents statistics for admin dashboard
type AdminStats struct {
	TotalRevenue  float64
	TotalOrders   int
	TotalProducts int
	TotalUsers    int
	RecentOrders  []Order
}
