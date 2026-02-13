package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// handleCart displays the shopping cart
func handleCart(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get or create cart for user
	cartID, err := getOrCreateCart(session.UserID)
	if err != nil {
		http.Error(w, "Error loading cart", http.StatusInternalServerError)
		log.Println("Cart error:", err)
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

	// Get cart item count
	cartCount := getCartItemCount(session.UserID)

	// Load template
	tmpl, err := template.New("cart.html").Funcs(templateFuncs).ParseFiles("templates/cart.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":    session.UserName,
		"Items":       items,
		"TotalAmount": totalAmount,
		"CartCount":   cartCount,
	}

	tmpl.Execute(w, data)
}

// handleAddToCart adds a product to the cart
func handleAddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := getSession(r)

	productID, _ := strconv.Atoi(r.FormValue("product_id"))
	size := r.FormValue("size")
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	if quantity == 0 {
		quantity = 1
	}

	// Get or create cart
	cartID, err := getOrCreateCart(session.UserID)
	if err != nil {
		http.Error(w, "Error accessing cart", http.StatusInternalServerError)
		log.Println("Cart error:", err)
		return
	}

	// Check if item already exists in cart
	var existingID int
	var existingQuantity int
	err = db.QueryRow(`
		SELECT id, quantity 
		FROM cart_items 
		WHERE cart_id = $1 AND product_id = $2 AND size = $3
	`, cartID, productID, size).Scan(&existingID, &existingQuantity)

	if err == sql.ErrNoRows {
		// Insert new item
		_, err = db.Exec(`
			INSERT INTO cart_items (cart_id, product_id, size, quantity)
			VALUES ($1, $2, $3, $4)
		`, cartID, productID, size, quantity)
	} else if err == nil {
		// Update existing item
		_, err = db.Exec(`
			UPDATE cart_items 
			SET quantity = quantity + $1
			WHERE id = $2
		`, quantity, existingID)
	}

	if err != nil {
		http.Error(w, "Error adding to cart", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Redirect to cart page
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

// handleRemoveFromCart removes an item from the cart
func handleRemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemID, _ := strconv.Atoi(r.FormValue("item_id"))

	_, err := db.Exec("DELETE FROM cart_items WHERE id = $1", itemID)
	if err != nil {
		http.Error(w, "Error removing item", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

// handleUpdateCart updates cart item quantity
func handleUpdateCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itemID, _ := strconv.Atoi(r.FormValue("item_id"))
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	if quantity < 1 {
		quantity = 1
	}

	_, err := db.Exec("UPDATE cart_items SET quantity = $1 WHERE id = $2", quantity, itemID)
	if err != nil {
		log.Println("Database error:", err)
	}

	// Return JSON response for AJAX
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// getOrCreateCart gets existing cart or creates a new one for the user
func getOrCreateCart(userID int) (int, error) {
	var cartID int

	// Try to get existing cart
	err := db.QueryRow("SELECT id FROM carts WHERE user_id = $1", userID).Scan(&cartID)

	if err == sql.ErrNoRows {
		// Create new cart
		err = db.QueryRow(
			"INSERT INTO carts (user_id) VALUES ($1) RETURNING id",
			userID,
		).Scan(&cartID)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}

	return cartID, nil
}
