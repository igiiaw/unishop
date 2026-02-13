package main

import (
	"html/template"
	"log"
	"net/http"
)

// handlePlaceOrder processes the order and clears the cart
func handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := getSession(r)

	// Get cart
	var cartID int
	err := db.QueryRow("SELECT id FROM carts WHERE user_id = $1", session.UserID).Scan(&cartID)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	// Get cart items
	rows, err := db.Query(`
		SELECT ci.product_id, ci.size, ci.quantity, p.name, p.price
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
	`, cartID)
	if err != nil {
		http.Error(w, "Error loading cart", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}
	defer rows.Close()

	// Calculate total and collect items
	type OrderItemData struct {
		ProductID   int
		ProductName string
		Size        string
		Quantity    int
		Price       float64
	}

	var items []OrderItemData
	var totalAmount float64

	for rows.Next() {
		var item OrderItemData
		err := rows.Scan(&item.ProductID, &item.Size, &item.Quantity, &item.ProductName, &item.Price)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		items = append(items, item)
		totalAmount += item.Price * float64(item.Quantity)
	}

	if len(items) == 0 {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error processing order", http.StatusInternalServerError)
		log.Println("Transaction error:", err)
		return
	}
	defer tx.Rollback()

	// Create order
	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, total_amount, status)
		VALUES ($1, $2, 'Completed')
		RETURNING id
	`, session.UserID, totalAmount).Scan(&orderID)
	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Insert order items
	for _, item := range items {
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, product_id, product_name, size, quantity, price)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, orderID, item.ProductID, item.ProductName, item.Size, item.Quantity, item.Price)
		if err != nil {
			http.Error(w, "Error saving order items", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}
	}

	// Clear cart
	_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		http.Error(w, "Error clearing cart", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Create notification for order
	_, err = tx.Exec(`
		INSERT INTO notifications (user_id, notification_type, message)
		VALUES ($1, 'order', $2)
	`, session.UserID, "Your order has been successfully placed and completed!")
	if err != nil {
		log.Println("Error creating notification:", err)
		// Continue anyway - notification is not critical
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Error completing order", http.StatusInternalServerError)
		log.Println("Commit error:", err)
		return
	}

	// Redirect to orders page
	http.Redirect(w, r, "/orders", http.StatusSeeOther)
}

// handleOrders displays order history
func handleOrders(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Check for success message
	successMsg := ""
	if r.URL.Query().Get("success") == "true" {
		successMsg = "Order placed successfully! Payment confirmed."
	}

	// Get all orders for user
	rows, err := db.Query(`
		SELECT id, total_amount, status, delivery_address, phone_number, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, session.UserID)
	if err != nil {
		http.Error(w, "Error loading orders", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.TotalAmount, &order.Status,
			&order.DeliveryAddress, &order.PhoneNumber, &order.CreatedAt)
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
		var itemCount int
		for itemRows.Next() {
			var item OrderItem
			err := itemRows.Scan(&item.ProductName, &item.Size, &item.Quantity, &item.Price)
			if err != nil {
				log.Println("Scan error:", err)
				continue
			}
			items = append(items, item)
			itemCount += item.Quantity
		}
		itemRows.Close()

		order.Items = items
		orders = append(orders, order)
	}

	// Get cart item count
	cartCount := getCartItemCount(session.UserID)

	// Load template
	tmpl, err := template.New("orders.html").Funcs(templateFuncs).ParseFiles("templates/orders.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":   session.UserName,
		"Orders":     orders,
		"CartCount":  cartCount,
		"SuccessMsg": successMsg,
	}

	tmpl.Execute(w, data)
}
