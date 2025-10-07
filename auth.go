package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
)

// Simple in-memory session store
var activeSessions = make(map[string]bool)

// Configuration - you can change these
const (
	USERNAME = "admin"   // Change this to your preferred username
	PASSWORD = "kpop123" // Change this to your preferred password
)

// Template data structure
type TemplateData struct {
	Cards           []PhotoCard `json:"cards"`
	IsAuthenticated bool   `json:"is_authenticated"`
}

// Generate secure session ID
func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Check if user is authenticated
func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false
	}
	return activeSessions[cookie.Value]
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Show login form
		tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>✨ Login to PC Collection ✨</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <style>
        .girly-gradient { background: linear-gradient(135deg, #ff6b9d, #c44569, #f8b500, #ff9ff3); }
        .cute-shadow { box-shadow: 0 10px 25px rgba(255, 107, 157, 0.3); }
        body { background: linear-gradient(135deg, #ffeef8, #e8f4ff, #fff0f8); min-height: 100vh; }
    </style>
</head>
<body class="flex items-center justify-center min-h-screen">
    <div class="bg-white bg-opacity-90 backdrop-blur-sm rounded-3xl cute-shadow p-8 w-full max-w-md">
        <div class="text-center mb-8">
            <h1 class="text-4xl font-bold mb-4">
                <span class="girly-gradient bg-clip-text text-transparent">
                    ✨ Login ✨
                </span>
            </h1>
            <p class="text-pink-600">Access your photocard collection</p>
        </div>
        
        <form method="POST" class="space-y-6">
            <div>
                <label class="block text-pink-600 font-bold mb-2">
                    <i class="fas fa-user mr-2"></i>Username
                </label>
                <input type="text" name="username" required 
                       class="w-full p-4 border-2 border-pink-200 rounded-2xl focus:outline-none focus:ring-2 focus:ring-pink-400 focus:border-pink-400">
            </div>
            
            <div>
                <label class="block text-pink-600 font-bold mb-2">
                    <i class="fas fa-lock mr-2"></i>Password
                </label>
                <input type="password" name="password" required 
                       class="w-full p-4 border-2 border-pink-200 rounded-2xl focus:outline-none focus:ring-2 focus:ring-pink-400 focus:border-pink-400">
            </div>
            
            <button type="submit" 
                    class="w-full girly-gradient text-white py-4 rounded-2xl font-bold text-lg transition-all duration-200 transform hover:scale-105 cute-shadow">
                <i class="fas fa-sign-in-alt mr-2"></i>Login ✨
            </button>
        </form>
        
        <div class="text-center mt-6">
            <a href="/" class="text-pink-500 hover:text-pink-700">
                <i class="fas fa-arrow-left mr-1"></i>Back to Collection
            </a>
        </div>
    </div>
</body>
</html>`

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, tmpl)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == USERNAME && password == PASSWORD {
			// Create session
			sessionID := generateSessionID()
			activeSessions[sessionID] = true

			// Set session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Secure:   false, // Set to true in production with HTTPS
				MaxAge:   86400, // 24 hours
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login?error=1", http.StatusSeeOther)
		}
	}
}

// Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		delete(activeSessions, cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Authentication middleware
func requireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}