package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

const sessionCookieName = "session"

// createSession creates a new session for the user
func createSession(w http.ResponseWriter, userID int, userName, email, role string) error {
	session := Session{
		UserID:   userID,
		UserName: userName,
		Email:    email,
		Role:     role,
	}

	// Encode session to JSON
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Encode to base64
	sessionValue := base64.StdEncoding.EncodeToString(sessionJSON)

	// Create cookie
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionValue,
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	return nil
}

// getSession retrieves the current session from request
func getSession(r *http.Request) (Session, error) {
	var session Session

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return session, err
	}

	// Decode base64
	sessionJSON, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return session, err
	}

	// Decode JSON
	err = json.Unmarshal(sessionJSON, &session)
	if err != nil {
		return session, err
	}

	return session, nil
}

// clearSession removes the session cookie
func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}
