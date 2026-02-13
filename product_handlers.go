package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Template helper functions
var templateFuncs = template.FuncMap{
	"multiply": func(a, b interface{}) float64 {
		var fa, fb float64
		switch v := a.(type) {
		case float64:
			fa = v
		case int:
			fa = float64(v)
		}
		switch v := b.(type) {
		case float64:
			fb = v
		case int:
			fb = float64(v)
		}
		return fa * fb
	},
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"iterate": func(n int) []int {
		result := make([]int, n)
		for i := range result {
			result[i] = i + 1
		}
		return result
	},
	"slice": func(s string, start int) string {
		if start >= len(s) {
			return ""
		}
		return s[start:]
	},
	"len": func(s string) int {
		return len(s)
	},
}

// handleCatalog displays the product catalog (8 products)
func handleCatalog(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get all products from database
	rows, err := db.Query(`
		SELECT id, name, description, price, image_url, 
		       COALESCE(highlight_text, '') as highlight_text,
		       COALESCE(highlight_enabled, false) as highlight_enabled
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
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL,
			&p.HighlightText, &p.HighlightEnabled)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		products = append(products, p)
	}

	// Get cart item count for badge
	cartCount := getCartItemCount(session.UserID)

	// Load template
	tmpl, err := template.New("catalog.html").Funcs(templateFuncs).ParseFiles("templates/catalog.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":  session.UserName,
		"Products":  products,
		"CartCount": cartCount,
	}

	tmpl.Execute(w, data)
}

// handleProduct displays product detail page
func handleProduct(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Extract product ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/product/")
	productID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get product from database
	var product Product
	err = db.QueryRow(`
		SELECT id, name, description, price, image_url,
		       COALESCE(highlight_text, '') as highlight_text,
		       COALESCE(highlight_enabled, false) as highlight_enabled
		FROM products 
		WHERE id = $1
	`, productID).Scan(&product.ID, &product.Name, &product.Description, &product.Price,
		&product.ImageURL, &product.HighlightText, &product.HighlightEnabled)

	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		log.Println("Database error:", err)
		return
	}

	// Get cart item count for badge
	cartCount := getCartItemCount(session.UserID)

	// Load template
	tmpl, err := template.New("product.html").Funcs(templateFuncs).ParseFiles("templates/product.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":  session.UserName,
		"Product":   product,
		"CartCount": cartCount,
	}

	tmpl.Execute(w, data)
}

// getCartItemCount returns the total number of items in user's cart
func getCartItemCount(userID int) int {
	var count int
	err := db.QueryRow(`
		SELECT COALESCE(SUM(ci.quantity), 0)
		FROM carts c
		LEFT JOIN cart_items ci ON c.id = ci.cart_id
		WHERE c.user_id = $1
	`, userID).Scan(&count)

	if err != nil {
		log.Println("Error getting cart count:", err)
		return 0
	}
	return count
}
