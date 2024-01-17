package handlers

import (
	"forum/pkg/auth"
	"forum/pkg/db"
	"net/http"
	"strconv"
)

func HandleCreatePostPage(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := db.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		title := r.Form.Get("title")
		content := r.Form.Get("content")
		categoriesSelected := r.Form["category"] // Updated to get multiple categories
		img := ""
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}
		if content == "" || title == "" || len(categoriesSelected) == 0 {
			session.Values["notification"] = "Invalid form submission"
			session.Save(r, w)
			http.Redirect(w, r, "/create_post", http.StatusSeeOther)
			return
		}

		if r.Form.Get("img") != "" {
			img = r.Form.Get("img")
		}

		// Create a post for each selected category
		for _, category := range categoriesSelected {
			categoryID, _ := strconv.Atoi(category)
			result, err := db.AddPost(userID, categoryID, content, title, img)
			if err != nil || result == nil {
				session.Values["notification"] = "Failed to create post"
				break
			}
			session.Values["notification"] = "Post created successfully"
		}

		session.Save(r, w)
		http.Redirect(w, r, "/create_post", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title":        "Sample create post Title",
		"SessionData":  session.Values,
		"IsCreatePost": true,
		"Categories":   categories,
	}

	renderTemplate(w, "create_post.html", data)
}
