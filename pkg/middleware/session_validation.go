package middleware

import (
	"forum/pkg/auth"
	"net/http"
)

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add authentication code here
		session, _ := auth.Store.Get(r, "user-session")
		if session.Values["authenticated"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		handler(w, r)
	}
}
