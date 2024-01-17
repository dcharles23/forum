package handlers

import (
	"forum/pkg/auth"
	"forum/pkg/db"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	// Set the response content type
	w.Header().Set("Content-Type", "text/html")

	// Parse header, footer, and page templates
	tmplFiles := []string{
		"templates/layout.html", // Ensure "layout.html" is executed first
		"templates/components/head.html",
		"templates/components/header.html",
		"templates/components/footer.html",
		"templates/pages/" + tmpl,
	}

	// Parse the templates
	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the templates with data
	err = t.ExecuteTemplate(w, "layout", data) // Execute "layout" template first
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderLogRegTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	// Set the response content type
	w.Header().Set("Content-Type", "text/html")

	// Parse header, footer, and page templates
	tmplFiles := []string{
		"templates/layout_logreg.html", // Ensure "layout.html" is executed first
		"templates/components/head_logreg.html",
		"templates/components/header_logreg.html",
		"templates/pages/" + tmpl,
	}

	// Parse the templates
	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the templates with data
	err = t.ExecuteTemplate(w, "layout_logreg", data) // Execute "layout" template first
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Rin add layout_profile.html head_profile.html  header_profile.html
func renderProfileTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	// Set the response content type
	w.Header().Set("Content-Type", "text/html")

	// Parse header, footer, and page templates
	tmplFiles := []string{
		"templates/layout_profile_own.html", // Ensure "layout.html" is executed first
		"templates/components/head_profile.html",
		"templates/components/header_profile.html",
		"templates/pages/" + tmpl,
	}

	// Parse the templates
	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the templates with data
	err = t.ExecuteTemplate(w, "layout_profile_own", data) // Execute "layout" template first
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func extractID(path string) (int, error) {
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Return an error if the conversion fails
		return 0, err
	}

	return id, nil
}

// HandleLikePost handles the liking of a post.
// HandleLikePost handles the liking of a post.
// HandleLikePost handles the liking of a post.
func HandleLikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	session, err := auth.Store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(int) // Assuming user_id is stored as an int

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if the user has previously liked the post
	hasLiked, err := db.HasUserLikedPost(userID, postID)
	if err != nil {
		http.Error(w, "Failed to check previous like status", http.StatusInternalServerError)
		return
	}

	// Remove the like if the user had previously liked the post
	if hasLiked {
		_, err = db.RemoveLike(userID, postID)
		if err != nil {
			http.Error(w, "Failed to remove the previous like", http.StatusInternalServerError)
			return
		}

		// Redirect back to the main page or the post's page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user has previously disliked the post
	hasDisliked, err := db.HasUserDislikedPost(userID, postID)
	if err != nil {
		http.Error(w, "Failed to check previous dislike status", http.StatusInternalServerError)
		return
	}

	// Remove the dislike if the user had previously disliked the post
	if hasDisliked {
		_, err = db.RemoveDislike(userID, postID)
		if err != nil {
			http.Error(w, "Failed to remove the previous dislike", http.StatusInternalServerError)
			return
		}
	}

	// Add logic to handle liking the post (e.g., update the database for like)
	// For simplicity, let's assume you have a function in the db package to handle likes
	_, err = db.AddLike(userID, postID)
	if err != nil {
		http.Error(w, "Failed to like the post", http.StatusInternalServerError)
		return
	}

	// Redirect back to the main page or the post's page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleDislikePost handles the disliking of a post.
// HandleDislikePost handles the disliking of a post.
func HandleDislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	session, err := auth.Store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(int) // Assuming user_id is stored as an int

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Check if the user has previously disliked the post
	hasDisliked, err := db.HasUserDislikedPost(userID, postID)
	if err != nil {
		http.Error(w, "Failed to check previous dislike status", http.StatusInternalServerError)
		return
	}

	// Remove the dislike if the user had previously disliked the post
	if hasDisliked {
		_, err = db.RemoveDislike(userID, postID)
		if err != nil {
			http.Error(w, "Failed to remove the previous dislike", http.StatusInternalServerError)
			return
		}

		// Redirect back to the main page or the post's page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user has previously liked the post
	hasLiked, err := db.HasUserLikedPost(userID, postID)
	if err != nil {
		http.Error(w, "Failed to check previous like status", http.StatusInternalServerError)
		return
	}

	// Remove the like if the user had previously liked the post
	if hasLiked {
		_, err = db.RemoveLike(userID, postID)
		if err != nil {
			http.Error(w, "Failed to remove the previous like", http.StatusInternalServerError)
			return
		}
	}

	// Add logic to handle disliking the post (e.g., update the database for dislike)
	// For simplicity, let's assume you have a function in the db package to handle dislikes
	_, err = db.AddDislike(userID, postID)
	if err != nil {
		http.Error(w, "Failed to dislike the post", http.StatusInternalServerError)
		return
	}

	// Redirect back to the main page or the post's page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleSubmitComment handles the submission of a comment.
func HandleSubmitComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	session, err := auth.Store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(int) // Assuming user_id is stored as an int

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	comment := r.FormValue("comment")
	if comment == "" {
		http.Error(w, "Invalid comment", http.StatusBadRequest)
		return

	}

	// Add logic to handle submitting the comment (e.g., update the database)
	// For simplicity, let's assume you have a function in the db package to handle comments
	_, err = db.AddComment(userID, postID, comment)
	if err != nil {
		http.Error(w, "Failed to submit the comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the main page or the post's page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleLikeComment handles the liking of a comment.
func HandleLikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	session, err := auth.Store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(int) // Assuming user_id is stored as an int

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if the user has previously liked the comment
	hasLiked, err := db.HasUserLikedComment(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to check previous like status", http.StatusInternalServerError)
		return
	}

	// Remove the like if the user had previously liked the comment
	if hasLiked {
		_, err = db.RemoveCommentLike(userID, commentID)
		if err != nil {
			http.Error(w, "Failed to remove the previous like", http.StatusInternalServerError)
			return
		}

		// Redirect back to the main page or the comment's page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user has previously disliked the comment
	hasDisliked, err := db.HasUserDislikedComment(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to check previous dislike status", http.StatusInternalServerError)
		return
	}

	// Remove the dislike if the user had previously disliked the comment
	if hasDisliked {
		_, err = db.RemoveCommentDislike(userID, commentID)
		if err != nil {
			http.Error(w, "Failed to remove the previous dislike", http.StatusInternalServerError)
			return
		}
	}

	// Add logic to handle liking the comment (e.g., update the database for like)
	_, err = db.AddCommentLike(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to like the comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the main page or the comment's page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleDislikeComment handles the disliking of a comment.
func HandleDislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	session, err := auth.Store.Get(r, "user-session")
	if err != nil || session.Values["user_id"] == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := session.Values["user_id"].(int) // Assuming user_id is stored as an int

	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if the user has previously disliked the comment
	hasDisliked, err := db.HasUserDislikedComment(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to check previous dislike status", http.StatusInternalServerError)
		return
	}

	// Remove the dislike if the user had previously disliked the comment
	if hasDisliked {
		_, err = db.RemoveCommentDislike(userID, commentID)
		if err != nil {
			http.Error(w, "Failed to remove the previous dislike", http.StatusInternalServerError)
			return
		}

		// Redirect back to the main page or the comment's page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user has previously liked the comment
	hasLiked, err := db.HasUserLikedComment(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to check previous like status", http.StatusInternalServerError)
		return
	}

	// Remove the like if the user had previously liked the comment
	if hasLiked {
		_, err = db.RemoveCommentLike(userID, commentID)
		if err != nil {
			http.Error(w, "Failed to remove the previous like", http.StatusInternalServerError)
			return
		}
	}

	// Add logic to handle disliking the comment (e.g., update the database for dislike)
	_, err = db.AddCommentDislike(userID, commentID)
	if err != nil {
		http.Error(w, "Failed to dislike the comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the main page or the comment's page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
