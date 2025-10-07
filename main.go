/*
K-pop Photocard Collection Web Application
=========================================

This Go web server provides a complete solution for managing a K-pop photocard collection.
It features a RESTful API, file upload handling, and a beautiful web interface.

Features:
- Upload and store photocard images
- Automatic tag generation based on group, album, and member
- Search and filter functionality
- Edit and delete existing photocards
- JSONL database storage for simplicity
- Responsive web interface with "girly pop" aesthetic
- User authentication system

Author: Your Name
Created: October 2025
*/

package main

import (
	"bufio"         // For reading files line by line
	"crypto/rand"   // For secure session generation
	"encoding/hex"  // For hex encoding
	"encoding/json" // For JSON marshaling/unmarshaling
	"fmt"           // For string formatting
	"html/template" // For HTML template rendering
	"io"            // For file copying operations
	"log"           // For logging server events
	"net/http"      // For HTTP server functionality
	"os"            // For file system operations
	"path/filepath" // For file path manipulation
	"strconv"       // For string to integer conversions
	"strings"       // For string manipulation
	"time"          // For timestamp generation
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
	IsAuthenticated bool        `json:"is_authenticated"`
}

/*
PhotoCard represents a single K-pop photocard in our collection.
Each card contains metadata about the idol, album, and associated image.

The JSON tags ensure proper serialization when sending data to the frontend.
*/
type PhotoCard struct {
	Group  string   `json:"group"`  // K-pop group name (e.g., "ATEEZ", "NewJeans")
	Album  string   `json:"album"`  // Album name (e.g., "Zero : Fever Pt.2")
	Member string   `json:"member"` // Member/idol name (e.g., "Mingi", "Hanni")
	Copies int      `json:"copies"` // Number of copies owned (default: 1)
	Image  string   `json:"image"`  // Path to uploaded image file
	Tags   []string `json:"tags"`   // Auto-generated hashtags for searchability
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

/*
loadCards reads all photocard data from the JSONL file.
JSONL (JSON Lines) format stores one JSON object per line, making it easy to
append new cards without loading the entire file into memory.

Parameters:
- filename: path to the JSONL file (typically "cards.jsonl")

Returns:
- []PhotoCard: slice of all photocards in the collection
- error: any file reading or JSON parsing errors
*/
func loadCards(filename string) ([]PhotoCard, error) {
	// Open the JSONL file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Ensure file is closed when function exits

	var cards []PhotoCard
	scanner := bufio.NewScanner(file)

	// Read file line by line
	for scanner.Scan() {
		var card PhotoCard
		// Parse JSON from current line
		if err := json.Unmarshal(scanner.Bytes(), &card); err != nil {
			// Log parsing errors but continue processing other lines
			log.Printf("error parsing line: %v", err)
			continue
		}
		cards = append(cards, card)
	}

	// Return any scanning errors
	return cards, scanner.Err()
}

/*
saveCard appends a new photocard to the JSONL file.
Using append mode ensures we don't overwrite existing data.

Parameters:
- card: PhotoCard struct to save to the database

Returns:
- error: any file writing or JSON marshaling errors
*/
func saveCard(card PhotoCard) error {
	// Open file in append mode, create if it doesn't exist
	file, err := os.OpenFile("cards.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert PhotoCard struct to JSON
	data, err := json.Marshal(card)
	if err != nil {
		return err
	}

	// Write JSON data followed by newline (JSONL format)
	_, err = file.Write(append(data, '\n'))
	return err
}

/*
uploadHandler processes photocard upload requests from the web form.
This handles multipart form data including file uploads and metadata.

The handler performs several tasks:
1. Validates the HTTP method (must be POST)
2. Parses multipart form data (images + text fields)
3. Processes and validates form inputs
4. Generates automatic tags based on group/album/member
5. Saves uploaded image file with unique filename
6. Creates and saves PhotoCard record
7. Redirects user back to main page

Parameters:
- w: HTTP response writer for sending response to client
- r: HTTP request containing form data and uploaded file
*/
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests for security
	if r.Method != "POST" {
		// Redirect GET requests to home page where upload form is located
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse multipart form data with 10MB size limit
	// This prevents users from uploading excessively large files
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", 400)
		return
	}

	// Extract form field values
	group := r.FormValue("group")      // K-pop group name
	album := r.FormValue("album")      // Album name
	member := r.FormValue("member")    // Member name
	copiesStr := r.FormValue("copies") // Number of copies (as string)

	// Convert copies string to integer, default to 1 if invalid
	copies, err := strconv.Atoi(copiesStr)
	if err != nil {
		copies = 1
	}

	// Auto-generate hashtags from form data for better searchability
	// Tags help users filter and find specific cards quickly
	var tags []string
	if group != "" {
		// Remove spaces from group name for hashtag format (#ATEEZ, #NewJeans)
		tags = append(tags, "#"+strings.ReplaceAll(group, " ", ""))
	}
	if album != "" {
		// Remove spaces from album name for hashtag format (#ZeroFeverPt2)
		tags = append(tags, "#"+strings.ReplaceAll(album, " ", ""))
	}
	if member != "" {
		// Remove spaces from member name for hashtag format (#Mingi, #Hanni)
		tags = append(tags, "#"+strings.ReplaceAll(member, " ", ""))
	}

	// Process uploaded image file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", 400)
		return
	}
	defer file.Close()

	// Generate unique filename using current timestamp to prevent conflicts
	// Format: {timestamp}_{original_filename}
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)
	filepath := filepath.Join("static", filename) // Save to static directory

	// Create destination file for uploaded image
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Error saving file", 500)
		return
	}
	defer dst.Close()

	// Copy uploaded file data to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", 500)
		return
	}

	// Create PhotoCard struct with all collected data
	card := PhotoCard{
		Group:  group,
		Album:  album,
		Member: member,
		Copies: copies,
		Image:  "/static/" + filename, // Web-accessible path with leading slash
		Tags:   tags,
	}

	// Save photocard data to JSONL database
	err = saveCard(card)
	if err != nil {
		http.Error(w, "Error saving card data", 500)
		return
	}

	// Redirect user back to main page to see their new photocard
	// Using 303 status code for proper POST-redirect-GET pattern
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/*
getUniqueGroups extracts all unique K-pop group names from the collection.
This is used to populate the auto-complete functionality when users upload new cards.

Returns:
- []string: slice of unique group names
- error: any database loading errors
*/
func getUniqueGroups() ([]string, error) {
	// Load all cards from database
	cards, err := loadCards("cards.jsonl")
	if err != nil {
		return nil, err
	}

	// Use map to track unique group names (map keys are inherently unique)
	groupSet := make(map[string]bool)
	for _, card := range cards {
		groupSet[card.Group] = true
	}

	// Convert map keys to slice for JSON response
	var groups []string
	for group := range groupSet {
		groups = append(groups, group)
	}
	return groups, nil
}

/*
deleteCard removes a photocard from the collection by index position.
Since we're using JSONL format, this requires rewriting the entire file.

Parameters:
- index: zero-based position of card to delete

Returns:
- error: any database operation errors
*/
func deleteCard(index int) error {
	// Load current collection
	cards, err := loadCards("cards.jsonl")
	if err != nil {
		return err
	}

	// Validate index bounds to prevent out-of-range errors
	if index < 0 || index >= len(cards) {
		return fmt.Errorf("invalid index")
	}

	// Remove card at specified index using slice operations
	// cards[:index] = elements before index
	// cards[index+1:] = elements after index
	cards = append(cards[:index], cards[index+1:]...)

	// Rewrite entire JSONL file with remaining cards
	file, err := os.Create("cards.jsonl")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each remaining card as a JSON line
	for _, card := range cards {
		data, err := json.Marshal(card)
		if err != nil {
			continue // Skip cards that can't be marshaled
		}
		file.Write(append(data, '\n'))
	}
	return nil
}

/*
updateCard modifies an existing photocard's data by index position.
Like deletion, this requires rewriting the entire JSONL file.

Parameters:
- index: zero-based position of card to update
- updatedCard: new PhotoCard data to replace existing card

Returns:
- error: any database operation errors
*/
func updateCard(index int, updatedCard PhotoCard) error {
	// Load current collection
	cards, err := loadCards("cards.jsonl")
	if err != nil {
		return err
	}

	// Validate index bounds
	if index < 0 || index >= len(cards) {
		return fmt.Errorf("invalid index")
	}

	// Replace card at specified index with new data
	cards[index] = updatedCard

	// Rewrite entire JSONL file with updated collection
	file, err := os.Create("cards.jsonl")
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each card (including updated one) as a JSON line
	for _, card := range cards {
		data, err := json.Marshal(card)
		if err != nil {
			continue // Skip cards that can't be marshaled
		}
		file.Write(append(data, '\n'))
	}
	return nil
}

/*
apiHandler routes API requests to appropriate functions.
This provides RESTful endpoints for the frontend JavaScript to interact with.

Supported endpoints:
- GET /api/groups - Returns unique group names for auto-complete
- GET /api/cards - Returns all photocard data for filtering
- POST /api/delete - Deletes a specific photocard by index
- POST /api/update - Updates a specific photocard by index

Parameters:
- w: HTTP response writer
- r: HTTP request containing API endpoint and data
*/
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type for all API responses
	w.Header().Set("Content-Type", "application/json")

	// Route request based on URL path
	switch r.URL.Path {
	case "/api/groups":
		// Handle auto-complete group data request
		if r.Method == "GET" {
			groups, err := getUniqueGroups()
			if err != nil {
				http.Error(w, "Error getting groups", 500)
				return
			}
			// Encode groups slice as JSON response
			json.NewEncoder(w).Encode(groups)
		}

	case "/api/cards":
		// Handle request for all photocard data (used by filters)
		if r.Method == "GET" {
			cards, err := loadCards("cards.jsonl")
			if err != nil {
				http.Error(w, "Error loading cards", 500)
				return
			}
			// Encode cards slice as JSON response
			json.NewEncoder(w).Encode(cards)
		}

	case "/api/delete":
		// Handle photocard deletion request
		if r.Method == "POST" {
			// Parse request body to get card index
			var req struct {
				Index int `json:"index"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", 400)
				return
			}

			// Delete the specified card
			if err := deleteCard(req.Index); err != nil {
				http.Error(w, "Error deleting card", 500)
				return
			}
			w.WriteHeader(200) // Success response
		}

	case "/api/update":
		// Handle photocard update request
		if r.Method == "POST" {
			// Parse request body to get card index and new data
			var req struct {
				Index int       `json:"index"`
				Card  PhotoCard `json:"card"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request", 400)
				return
			}

			// Update the specified card with new data
			if err := updateCard(req.Index, req.Card); err != nil {
				http.Error(w, "Error updating card", 500)
				return
			}
			w.WriteHeader(200) // Success response
		}

	default:
		// Unknown API endpoint
		http.NotFound(w, r)
	}
}

/*
main function initializes and starts the web server.
This sets up all HTTP routes and begins listening for requests.
*/
func main() {
	// Parse HTML template file (panics if template has syntax errors)
	// template.Must ensures the server won't start with a broken template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Set up static file serving for uploaded images
	// This allows the web browser to access files in the static/ directory
	// URL pattern: /static/filename.jpg -> static/filename.jpg on disk
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Main page route - displays the photocard collection interface
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Load all photocards from database
		cards, err := loadCards("cards.jsonl")
		if err != nil {
			http.Error(w, "Error loading cards", 500)
			return
		}

		// Check if user is authenticated
		userAuthenticated := isAuthenticated(r)

		// Create template data with both cards and authentication status
		templateData := TemplateData{
			Cards:           cards,
			IsAuthenticated: userAuthenticated,
		}

		// Render HTML template with complete data
		err = tmpl.Execute(w, templateData)
		if err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Template error", 500)
			return
		}
	})

	// Photo upload route - handles new photocard submissions
	http.HandleFunc("/upload", requireAuth(uploadHandler))

	// API routes - provides JSON endpoints for frontend JavaScript
	http.HandleFunc("/api/", requireAuth(apiHandler))

	// Authentication routes
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Get port from environment or use default
	port := getPort()

	// Start HTTP server and log startup message
	log.Println("🌸 K-pop Photocard Collection Server Started! 🌸")
	log.Printf("📍 Server running at: http://localhost%s", port)
	log.Println("📁 Images stored in: ./static/")
	log.Println("💾 Database file: ./cards.jsonl")
	log.Println("✨ Ready to collect some precious photocards! ✨")

	// ListenAndServe blocks forever, serving HTTP requests
	// log.Fatal will print any server startup errors and exit
	log.Fatal(http.ListenAndServe(port, nil))
}
