package model

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	UserID        int64  `json:"user_id"`
	User          User   `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}
