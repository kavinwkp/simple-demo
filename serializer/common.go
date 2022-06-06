package serializer

import "github.com/RaymondCode/simple-demo/model"

// 基础序列化器
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User model.User `json:"user"`
}

type PublishResponse struct {
	Response
}

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}
