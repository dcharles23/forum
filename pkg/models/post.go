package models

type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	Img          string
	CategoryID   int
	LikeCount    int
	DislikeCount int
}
