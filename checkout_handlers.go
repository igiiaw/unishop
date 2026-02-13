package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// handleCheckout displays the checkout page with saved card info
func handleCheckout(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get cart items to verify cart is not empty
	cartID, err := getOrCreateCart(session.UserID)
	if err != nil {
		http.Error(w, "Error loading cart", http.StatusInternalServerError)
		return
	}

	// Get cart items with product details
	rows, err := db.Query(`
       SELECT ci.id, ci.product_id, ci.size, ci.quantity,
              p.name, p.price, p.image_url
       FROM cart_items ci
       JOIN products p ON ci.product_id = p.id
       WHERE ci.cart_id = $1
       ORDER BY ci.added_at DESC
    `, cartID)
	if err != nil {
		http.Error(w, "Error loading cart items", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}
	defer rows.Close()

	var items []CartItem
	var totalAmount float64

	for rows.Next() {
		var item CartItem
		item.Product = &Product{}

		err := rows.Scan(
			&item.ID, &item.ProductID, &item.Size, &item.Quantity,
			&item.Product.Name, &item.Product.Price, &item.Product.ImageURL,
		)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		item.CartID = cartID
		items = append(items, item)
		totalAmount += item.Product.Price * float64(item.Quantity)
	}

	// If cart is empty, redirect to cart page
	if len(items) == 0 {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	// Get saved card if exists
	var savedCard *SavedCard
	var card SavedCard
	err = db.QueryRow(`
       SELECT id, cardholder_name, card_number, expiry_month, expiry_year, cvv
       FROM saved_cards
       WHERE user_id = $1
    `, session.UserID).Scan(&card.ID, &card.CardholderName, &card.CardNumber,
		&card.ExpiryMonth, &card.ExpiryYear, &card.CVV)

	if err == nil {
		savedCard = &card
	} else if err != sql.ErrNoRows {
		log.Println("Error fetching saved card:", err)
	}

	// Get cart item count
	cartCount := getCartItemCount(session.UserID)

	// === НАЧАЛО ИСПРАВЛЕНИЯ ===
	// 1. Определяем необходимые функции для шаблона
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(price float64, qty int) float64 {
			return price * float64(qty)
		},
		"iterate": func(count int) []int {
			var i int
			var items []int
			for i = 1; i <= count; i++ {
				items = append(items, i)
			}
			return items
		},
	}

	// 2. Используем .Funcs(funcMap) при загрузке шаблона
	tmpl, err := template.New("checkout.html").Funcs(funcMap).ParseFiles("templates/checkout.html")
	// === КОНЕЦ ИСПРАВЛЕНИЯ ===

	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":    session.UserName,
		"Items":       items,
		"TotalAmount": totalAmount,
		"SavedCard":   savedCard,
		"CartCount":   cartCount,
	}

	tmpl.Execute(w, data)
}

// handleProcessCheckout processes the payment and creates order
func handleProcessCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := getSession(r)

	// Parse form data
	deliveryAddress := r.FormValue("delivery_address")
	phoneNumber := r.FormValue("phone_number")
	cardholderName := r.FormValue("cardholder_name")
	cardNumber := r.FormValue("card_number")
	expiryMonth := r.FormValue("expiry_month")
	expiryYear := r.FormValue("expiry_year")
	cvv := r.FormValue("cvv")
	saveCard := r.FormValue("save_card") == "on"

	// Validate inputs
	if deliveryAddress == "" || phoneNumber == "" {
		http.Error(w, "Delivery address and phone number are required", http.StatusBadRequest)
		return
	}

	if cardholderName == "" || cardNumber == "" || expiryMonth == "" || expiryYear == "" || cvv == "" {
		http.Error(w, "All card fields are required", http.StatusBadRequest)
		return
	}

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

	var orderItems []OrderItemData
	var totalAmount float64

	for rows.Next() {
		var item OrderItemData
		err := rows.Scan(&item.ProductID, &item.Size, &item.Quantity, &item.ProductName, &item.Price)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		orderItems = append(orderItems, item)
		totalAmount += item.Price * float64(item.Quantity)
	}

	if len(orderItems) == 0 {
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

	// Save card if requested
	if saveCard {
		// Check if card already exists
		var existingCardID int
		err = tx.QueryRow("SELECT id FROM saved_cards WHERE user_id = $1", session.UserID).Scan(&existingCardID)

		if err == sql.ErrNoRows {
			// Insert new card
			_, err = tx.Exec(`
				INSERT INTO saved_cards (user_id, cardholder_name, card_number, expiry_month, expiry_year, cvv)
				VALUES ($1, $2, $3, $4, $5, $6)
			`, session.UserID, cardholderName, cardNumber, expiryMonth, expiryYear, cvv)
		} else if err == nil {
			// Update existing card
			_, err = tx.Exec(`
				UPDATE saved_cards 
				SET cardholder_name = $1, card_number = $2, expiry_month = $3, expiry_year = $4, cvv = $5, updated_at = CURRENT_TIMESTAMP
				WHERE user_id = $6
			`, cardholderName, cardNumber, expiryMonth, expiryYear, cvv, session.UserID)
		}

		if err != nil {
			log.Println("Error saving card:", err)
			// Continue anyway - card saving is optional
		}
	}

	// Create order with delivery info
	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, total_amount, status, delivery_address, phone_number, payment_method)
		VALUES ($1, $2, 'PAID', $3, $4, 'CARD')
		RETURNING id
	`, session.UserID, totalAmount, deliveryAddress, phoneNumber).Scan(&orderID)
	if err != nil {
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Insert order items
	for _, item := range orderItems {
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
	`, session.UserID, "Your order #"+fmt.Sprint(orderID)+" has been successfully placed and paid! Total: $"+fmt.Sprintf("%.2f", totalAmount))
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
	http.Redirect(w, r, "/orders?success=true", http.StatusSeeOther)
}

// handleSaveCard saves or updates user's card
func handleSaveCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := getSession(r)

	cardholderName := r.FormValue("cardholder_name")
	cardNumber := r.FormValue("card_number")
	expiryMonth := r.FormValue("expiry_month")
	expiryYear := r.FormValue("expiry_year")
	cvv := r.FormValue("cvv")

	// Validate inputs
	if cardholderName == "" || cardNumber == "" || expiryMonth == "" || expiryYear == "" || cvv == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Check if card already exists
	var existingCardID int
	err := db.QueryRow("SELECT id FROM saved_cards WHERE user_id = $1", session.UserID).Scan(&existingCardID)

	if err == sql.ErrNoRows {
		// Insert new card
		_, err = db.Exec(`
			INSERT INTO saved_cards (user_id, cardholder_name, card_number, expiry_month, expiry_year, cvv)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, session.UserID, cardholderName, cardNumber, expiryMonth, expiryYear, cvv)
	} else if err == nil {
		// Update existing card
		_, err = db.Exec(`
			UPDATE saved_cards 
			SET cardholder_name = $1, card_number = $2, expiry_month = $3, expiry_year = $4, cvv = $5, updated_at = CURRENT_TIMESTAMP
			WHERE user_id = $6
		`, cardholderName, cardNumber, expiryMonth, expiryYear, cvv, session.UserID)
	}

	if err != nil {
		http.Error(w, "Error saving card", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Redirect back to checkout
	http.Redirect(w, r, "/checkout", http.StatusSeeOther)
}

// maskCardNumber masks card number for display (shows last 4 digits)
func maskCardNumber(cardNumber string) string {
	if len(cardNumber) < 4 {
		return "****"
	}
	return "•••• •••• •••• " + cardNumber[len(cardNumber)-4:]
}
