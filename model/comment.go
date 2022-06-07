package model

type Comment struct {
	Id      int64  `json:"id,omitempty"`
	UserID  int64  `json:"user_id"`
	VideoID int64  `json:"video_id"`
	Content string `json:"content,omitempty"`
}

type CommentResponse struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
