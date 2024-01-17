package handlers

import "net/http"

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Login Page",
		"Content": "Please log in to access the forum.",
	}
	renderLogRegTemplate(w, "login.html", data)
}
