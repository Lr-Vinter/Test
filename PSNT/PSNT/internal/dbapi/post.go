package dbapi

type Post struct {
	PostID     int
	OwnerID    int
	CreatedAt  int64
	LikeNumber int
	Message    string
}

type Selection []Post
