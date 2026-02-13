package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// handleAdminDashboard displays admin dashboard with statistics
func handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get statistics
	var stats AdminStats

	// Total revenue
	db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders").Scan(&stats.TotalRevenue)

	// Total orders
	db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)

	// Total products
	db.QueryRow("SELECT COUNT(*) FROM products").Scan(&stats.TotalProducts)

	// Total users
	db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'USER'").Scan(&stats.TotalUsers)

	// Recent orders (last 10)
	rows, err := db.Query(`
		SELECT o.id, o.user_id, o.total_amount, o.status, o.delivery_address, 
		       o.phone_number, o.created_at, u.name
		FROM orders o
		JOIN users u ON o.user_id = u.id
		ORDER BY o.created_at DESC
		LIMIT 10
	`)
	if err != nil {
		log.Println("Error fetching orders:", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var order Order
			rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status,
				&order.DeliveryAddress, &order.PhoneNumber, &order.CreatedAt, &order.UserName)
			stats.RecentOrders = append(stats.RecentOrders, order)
		}
	}

	// Load template
	tmpl, err := template.ParseFiles("templates/admin_dashboard.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName": session.UserName,
		"Stats":    stats,
	}

	tmpl.Execute(w, data)
}

// handleAdminProducts displays product management page
func handleAdminProducts(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get all products
	rows, err := db.Query(`
		SELECT id, name, description, price, image_url, 
		       COALESCE(highlight_text, '') as highlight_text,
		       COALESCE(highlight_enabled, false) as highlight_enabled,
		       created_at, updated_at
		FROM products
		ORDER BY id
	`)
	if err != nil {
		http.Error(w, "Error loading products", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var updatedAt sql.NullTime
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL,
			&p.HighlightText, &p.HighlightEnabled, &p.CreatedAt, &updatedAt)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		products = append(products, p)
	}

	// Load template
	tmpl, err := template.ParseFiles("templates/admin_products.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName": session.UserName,
		"Products": products,
	}

	tmpl.Execute(w, data)
}

// handleAdminAddProduct handles adding a new product
func handleAdminAddProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		session, _ := getSession(r)

		tmpl, err := template.ParseFiles("templates/admin_product_form.html")
		if err != nil {
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		data := map[string]interface{}{
			"UserName": session.UserName,
			"Action":   "Add",
			"Product":  nil,
		}

		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		description := r.FormValue("description")
		priceStr := r.FormValue("price")
		imageURL := r.FormValue("image_url")
		highlightText := r.FormValue("highlight_text")
		highlightEnabled := r.FormValue("highlight_enabled") == "on"

		// Validate
		if name == "" || priceStr == "" {
			http.Error(w, "Name and price are required", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		if imageURL == "" {
			imageURL = "/static/images/placeholder.jpg"
		}

		// Insert product
		_, err = db.Exec(`
			INSERT INTO products (name, description, price, image_url, highlight_text, highlight_enabled)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, name, description, price, imageURL, highlightText, highlightEnabled)

		if err != nil {
			http.Error(w, "Error adding product", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}

		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleAdminEditProduct handles editing a product
func handleAdminEditProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	if productID == "" {
		http.Error(w, "Product ID required", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		session, _ := getSession(r)

		// Get product
		var product Product
		err := db.QueryRow(`
			SELECT id, name, description, price, image_url,
			       COALESCE(highlight_text, '') as highlight_text,
			       COALESCE(highlight_enabled, false) as highlight_enabled
			FROM products
			WHERE id = $1
		`, productID).Scan(&product.ID, &product.Name, &product.Description, &product.Price,
			&product.ImageURL, &product.HighlightText, &product.HighlightEnabled)

		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles("templates/admin_product_form.html")
		if err != nil {
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		data := map[string]interface{}{
			"UserName": session.UserName,
			"Action":   "Edit",
			"Product":  product,
		}

		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		description := r.FormValue("description")
		priceStr := r.FormValue("price")
		imageURL := r.FormValue("image_url")
		highlightText := r.FormValue("highlight_text")
		highlightEnabled := r.FormValue("highlight_enabled") == "on"

		// Validate
		if name == "" || priceStr == "" {
			http.Error(w, "Name and price are required", http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		// Update product
		_, err = db.Exec(`
			UPDATE products 
			SET name = $1, description = $2, price = $3, image_url = $4, 
			    highlight_text = $5, highlight_enabled = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
		`, name, description, price, imageURL, highlightText, highlightEnabled, productID)

		if err != nil {
			http.Error(w, "Error updating product", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}

		http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleAdminDeleteProduct handles deleting a product
func handleAdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	productID := r.FormValue("id")
	if productID == "" {
		http.Error(w, "Product ID required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

// handleAdminOrders displays all orders for admin
func handleAdminOrders(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get all orders
	rows, err := db.Query(`
		SELECT o.id, o.user_id, o.total_amount, o.status, o.delivery_address, 
		       o.phone_number, o.payment_method, o.created_at, u.name
		FROM orders o
		JOIN users u ON o.user_id = u.id
		ORDER BY o.created_at DESC
	`)
	if err != nil {
		http.Error(w, "Error loading orders", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status,
			&order.DeliveryAddress, &order.PhoneNumber, &order.PaymentMethod,
			&order.CreatedAt, &order.UserName)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		// Get order items
		itemRows, err := db.Query(`
			SELECT product_name, size, quantity, price
			FROM order_items
			WHERE order_id = $1
		`, order.ID)
		if err != nil {
			log.Println("Error loading order items:", err)
			continue
		}

		var items []OrderItem
		for itemRows.Next() {
			var item OrderItem
			err := itemRows.Scan(&item.ProductName, &item.Size, &item.Quantity, &item.Price)
			if err != nil {
				log.Println("Scan error:", err)
				continue
			}
			items = append(items, item)
		}
		itemRows.Close()

		order.Items = items
		orders = append(orders, order)
	}

	// Load template
	tmpl, err := template.New("admin_orders.html").Funcs(templateFuncs).ParseFiles("templates/admin_orders.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName": session.UserName,
		"Orders":   orders,
	}

	tmpl.Execute(w, data)
}

// handleAdminNotifications displays notification management page
func handleAdminNotifications(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Check for success message
	successMsg := ""
	if r.URL.Query().Get("success") == "true" {
		successMsg = "Notification sent successfully to all users!"
	}

	// Get recent notifications sent by admin (global notifications)
	rows, err := db.Query(`
		SELECT id, notification_type, message, created_at
		FROM notifications
		WHERE is_global = true
		ORDER BY created_at DESC
		LIMIT 20
	`)

	var recentNotifications []Notification
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var n Notification
			rows.Scan(&n.ID, &n.NotificationType, &n.Message, &n.CreatedAt)
			n.IsGlobal = true
			recentNotifications = append(recentNotifications, n)
		}
	} else {
		log.Println("Error fetching notifications:", err)
	}

	// Load template
	tmpl, err := template.ParseFiles("templates/admin_notifications.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":            session.UserName,
		"SuccessMsg":          successMsg,
		"RecentNotifications": recentNotifications,
	}

	tmpl.Execute(w, data)
}

// handleAdminSendNotification processes sending a notification to all users
func handleAdminSendNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	message := r.FormValue("message")

	// Validate input
	if message == "" {
		http.Error(w, "Error: Message is required", http.StatusBadRequest)
		return
	}

	// Insert global notification
	_, err := db.Exec(`
		INSERT INTO notifications (user_id, notification_type, message, is_global)
		VALUES (NULL, 'promo', $1, TRUE)
	`, message)

	if err != nil {
		http.Error(w, "Error sending notification", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Redirect back with success message
	http.Redirect(w, r, "/admin/notifications?success=true", http.StatusSeeOther)
}
