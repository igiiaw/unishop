package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// handleHome redirects to login page if not authenticated, otherwise to catalog
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	session, err := getSession(r)
	if err != nil || session.UserID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

// handleRegister handles user registration
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Show registration form
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		data := map[string]interface{}{
			"Error": "",
		}
		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		// Process registration
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate input
		if name == "" || email == "" || password == "" {
			showRegisterError(w, "All fields are required")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error processing registration", http.StatusInternalServerError)
			log.Println("Bcrypt error:", err)
			return
		}

		// Insert user into database
		_, err = db.Exec(
			"INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)",
			name, email, string(hashedPassword),
		)
		if err != nil {
			showRegisterError(w, "Email already exists")
			log.Println("Database error:", err)
			return
		}

		// Redirect to login page
		http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleLogin handles user authentication
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// If already logged in, redirect to catalog
	session, _ := getSession(r)
	if session.UserID != 0 {
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		// Show login form
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Error loading page", http.StatusInternalServerError)
			log.Println("Template error:", err)
			return
		}

		registered := r.URL.Query().Get("registered") == "true"
		data := map[string]interface{}{
			"Error":      "",
			"Registered": registered,
		}
		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {
		// Process login
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate input
		if email == "" || password == "" {
			showLoginError(w, "Email and password are required")
			return
		}

		// Get user from database
		var user User
		err := db.QueryRow(
			"SELECT id, name, email, password_hash, role FROM users WHERE email = $1",
			email,
		).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

		if err == sql.ErrNoRows {
			showLoginError(w, "Invalid email or password")
			return
		}
		if err != nil {
			http.Error(w, "Error processing login", http.StatusInternalServerError)
			log.Println("Database error:", err)
			return
		}

		// Verify password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			showLoginError(w, "Invalid email or password")
			return
		}

		// Create session
		err = createSession(w, user.ID, user.Name, user.Email, user.Role)
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			log.Println("Session error:", err)
			return
		}

		// Redirect based on role
		if user.Role == "ADMIN" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/catalog", http.StatusSeeOther)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleLogout handles user logout
func handleLogout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// showRegisterError displays registration page with error message
func showRegisterError(w http.ResponseWriter, errorMsg string) {
	tmpl, _ := template.ParseFiles("templates/register.html")
	data := map[string]interface{}{
		"Error": errorMsg,
	}
	tmpl.Execute(w, data)
}

// showLoginError displays login page with error message
func showLoginError(w http.ResponseWriter, errorMsg string) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	data := map[string]interface{}{
		"Error":      errorMsg,
		"Registered": false,
	}
	tmpl.Execute(w, data)
}
