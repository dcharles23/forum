package handlers

import "net/http"

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Registration Page",
		"Content": "Create an account to join the forum.",
	}
	renderLogRegTemplate(w, "register.html", data)
}
