package main

import (
	"html/template"
	"log"
	"net/http"
)

// handleProfile displays user profile page
func handleProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get user data
	var user User
	err := db.QueryRow(`
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`, session.UserID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		http.Error(w, "Error loading profile", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}

	// Get order count
	var orderCount int
	db.QueryRow("SELECT COUNT(*) FROM orders WHERE user_id = $1", session.UserID).Scan(&orderCount)

	// Get cart item count
	cartCount := getCartItemCount(session.UserID)

	// Load template
	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":   session.UserName,
		"User":       user,
		"OrderCount": orderCount,
		"CartCount":  cartCount,
	}

	tmpl.Execute(w, data)
}

// handleNotifications displays user notifications
func handleNotifications(w http.ResponseWriter, r *http.Request) {
	session, _ := getSession(r)

	// Get notifications (both personal and global)
	rows, err := db.Query(`
		SELECT id, notification_type, message, is_global, created_at
		FROM notifications
		WHERE user_id = $1 OR is_global = true
		ORDER BY created_at DESC
		LIMIT 50
	`, session.UserID)
	if err != nil {
		http.Error(w, "Error loading notifications", http.StatusInternalServerError)
		log.Println("Database error:", err)
		return
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err := rows.Scan(&n.ID, &n.NotificationType, &n.Message, &n.IsGlobal, &n.CreatedAt)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}
		notifications = append(notifications, n)
	}

	// Get cart item count
	cartCount := getCartItemCount(session.UserID)

	// Load template
	tmpl, err := template.ParseFiles("templates/notifications.html")
	if err != nil {
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		log.Println("Template error:", err)
		return
	}

	data := map[string]interface{}{
		"UserName":      session.UserName,
		"Notifications": notifications,
		"CartCount":     cartCount,
	}

	tmpl.Execute(w, data)
}
