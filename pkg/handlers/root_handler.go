package handlers

import (
	"fmt"
	"forum/pkg/auth"
	"forum/pkg/db"
	"forum/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	// You can retrieve the post data from the database here
	session, err := auth.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Declare variables for the presence of specific query parameters
	var (
		categoryIDExists bool
		likesExist       bool
		postsByUser      bool
	)

	// Check if specific query parameters exist
	if _, ok := r.URL.Query()["category_id"]; ok {
		categoryIDExists = true
		fmt.Println("category_id exists")
	}

	if _, ok := r.URL.Query()["likes"]; ok {
		likesExist = true
		fmt.Println("likes exists")
	}

	if _, ok := r.URL.Query()["created"]; ok {
		postsByUser = true
		fmt.Println("comments exists")
	}

	// Declare a variable to hold posts
	var posts []models.Post

	// Check all errors here and handle them i didnt do it @@@@
	var allCategories []models.Category
	allCategories, _ = db.GetAllCategories()

	if categoryIDExists || likesExist || postsByUser {
		// If any of the query parameters exist, use filtered posts
		if session.Values["username"] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		userIDm, err := db.GetUserIDByUsername(session.Values["username"].(string)) // Get user ID
		if err != nil {
			// Handle the error appropriately, such as logging or returning an error response to the client
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// If any of the query parameters exist, use filtered posts
		intCategoryID := 0
		if categoryIDExists {
			categoryID := r.URL.Query().Get("category_id")
			var err error
			intCategoryID, err = strconv.Atoi(categoryID)
			if err != nil {
				renderTemplate(w, "404.html", nil)
				return
			}
			fmt.Println("category ID:", intCategoryID)
		}

		// Get filtered posts
		filteredPosts, err := db.GetPostsByFilters(intCategoryID, likesExist, postsByUser, userIDm)
		if err != nil {
			// Handle the error appropriately, such as logging or returning an error response to the client
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Use filtered posts
		posts = filteredPosts
	} else {
		// If no query parameters, get all posts
		var allPosts []models.Post
		allPosts, _ = db.GetAllPosts()
		posts = allPosts
	}

	likeCounts := make(map[int]int)
	// Loop through all posts and count the number of likes for each post
	for _, post := range posts {
		postID := post.ID
		likeCount, err := db.CountLikesByPost(postID)
		if err != nil {
			// Handle the error if necessary
			log.Println("Error counting likes for post", postID, ":", err)
			continue // Continue to the next iteration
		}
		likeCounts[postID] = likeCount
	}

	dislikeCounts := make(map[int]int)
	// Loop through all posts and count the number of likes for each post
	for _, post := range posts {
		postID := post.ID
		dislikeCount, err := db.CountDisLikesByPost(postID)
		if err != nil {
			// Handle the error if necessary
			log.Println("Error counting likes for post", postID, ":", err)
			continue // Continue to the next iteration
		}
		dislikeCounts[postID] = dislikeCount
	}
	// User is authenticated
	// Pass the session data to the template
	data := map[string]interface{}{
		"Title":       "Home Page",
		"SessionData": session.Values,
		"IsHomePage":  true,
		"Posts":       posts,
		"Categories":  allCategories,
		"Likes":       likeCounts,
		"Dislikes":    dislikeCounts,
	}
	renderTemplate(w, "index.html", data)
}
