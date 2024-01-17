package models

type Comment struct {
	ID           int
	UserID       int // Assuming user_id is stored as an int
	PostID       int
	Content      string
	LikeCount    int
	DislikeCount int
	// ENd of comment
}
