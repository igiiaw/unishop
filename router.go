package main

import (
	"net/http"
)

// NewRouter creates and configures the HTTP router
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes (no authentication required)
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/register", handleRegister)
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/logout", handleLogout)

	// Protected routes (authentication required)
	mux.HandleFunc("/catalog", requireAuth(handleCatalog))
	mux.HandleFunc("/product/", requireAuth(handleProduct))
	mux.HandleFunc("/cart", requireAuth(handleCart))
	mux.HandleFunc("/cart/add", requireAuth(handleAddToCart))
	mux.HandleFunc("/cart/remove", requireAuth(handleRemoveFromCart))
	mux.HandleFunc("/cart/update", requireAuth(handleUpdateCart))
	mux.HandleFunc("/checkout", requireAuth(handleCheckout))
	mux.HandleFunc("/checkout/process", requireAuth(handleProcessCheckout))
	mux.HandleFunc("/payment/card", requireAuth(handleSaveCard))
	mux.HandleFunc("/orders", requireAuth(handleOrders))
	mux.HandleFunc("/profile", requireAuth(handleProfile))
	mux.HandleFunc("/notifications", requireAuth(handleNotifications))

	// Admin routes (require admin role)
	mux.HandleFunc("/admin", requireAdmin(handleAdminDashboard))
	mux.HandleFunc("/admin/products", requireAdmin(handleAdminProducts))
	mux.HandleFunc("/admin/products/add", requireAdmin(handleAdminAddProduct))
	mux.HandleFunc("/admin/products/edit", requireAdmin(handleAdminEditProduct))
	mux.HandleFunc("/admin/products/delete", requireAdmin(handleAdminDeleteProduct))
	mux.HandleFunc("/admin/orders", requireAdmin(handleAdminOrders))
	mux.HandleFunc("/admin/notifications", requireAdmin(handleAdminNotifications))
	mux.HandleFunc("/admin/notifications/send", requireAdmin(handleAdminSendNotification))

	return mux
}

// requireAuth is a middleware that checks if user is authenticated
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		session, err := getSession(r)
		if err != nil || session.UserID == 0 {
			// Not authenticated - redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// User is authenticated - continue to handler
		next(w, r)
	}
}

// requireAdmin is a middleware that checks if user is an admin
func requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		session, err := getSession(r)
		if err != nil || session.UserID == 0 {
			// Not authenticated - redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Check if user is admin
		if session.Role != "ADMIN" {
			// Not admin - redirect to catalog with error
			http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
			return
		}
		// User is admin - continue to handler
		next(w, r)
	}
}
