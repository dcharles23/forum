package handlers

import "net/http"

// Handle when user is looking at their own profile
func HandleOwnProfile(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Profile Page",
		"Content": "This is a sample Profile.",
	}
	renderProfileTemplate(w, "user_profile_own.html", data)
}

// Handle when user is looking at other user's profile
func HandleOtherProfile(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Profile Page",
		"Content": "This is a sample Profile.",
	}
	renderProfileTemplate(w, "user_profile_other.html", data)
}
