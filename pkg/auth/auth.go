package auth

import (
	"database/sql"
	"fmt"
	"forum/pkg/db" // Replace with the actual path to your db package
	"log"
	"net/http"
	"unicode"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Store will hold all session data
var Store = sessions.NewCookieStore([]byte("super-secret-key"))

func init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
		Secure:   true, // Use Secure flag for production
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordSpace := false

	for _, char := range password {
		if unicode.IsSpace(int32(char)) {
			passwordSpace = true
		}
	}

	if db.CheckUserExists(email, username) {
		http.Error(w, "User Already Exists", http.StatusInternalServerError)
		return
	}

	if passwordSpace {
		http.Error(w, "Password cannot contain spaces!", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		return
	}

	_, err = db.AddUser(email, username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// LoginHandler handles the user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	fmt.Println(email, ":::", password)
	row, err := db.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	var (
		userID         int
		userEmail      string
		username       string
		storedPassword string
	)

	if err := row.Scan(&userID, &userEmail, &username, &storedPassword); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("Error scanning user data: %v", err)
		http.Error(w, "Login error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	sessionUUID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error generating session UUID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session, _ := Store.Get(r, "user-session")
	session.Values["uuid"] = sessionUUID.String()
	session.Values["user_id"] = userID
	session.Values["username"] = username
	session.Values["email"] = userEmail
	session.Values["authenticated"] = true
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler handles logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Set the "authenticated" value to false to indicate that the user is logged out.
	session.Values["authenticated"] = false

	// Delete the session values you want to clear (e.g., "uuid" and "authenticated").
	delete(session.Values, "uuid")
	delete(session.Values, "authenticated")
	delete(session.Values, "user_id")
	delete(session.Values, "username")
	delete(session.Values, "email")

	// Save the session to apply the changes.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	// Redirect the user to the login page or any other appropriate destination.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
