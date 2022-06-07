package service

import (
	"time"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/utils"
)

type CommentService struct {
	Token     string `form:"token" json:"token"`
	VideoID   int64  `form:"video_id" json:"video_id"`
	Content   string `form:"user_name" json:"user_name"`
	CommentID int64  `form:"comment_id" json:"comment_id"`
}

func (service *CommentService) Comment() serializer.CommentActionResponse {
	claim, _ := utils.ParseToken(service.Token)
	user_id := claim.UserId

	video_id := service.VideoID
	content := service.Content

	var user model.User
	model.DB.First(&user, user_id)

	var comment = model.Comment{
		UserID:  user_id,
		VideoID: video_id,
		Content: content,
	}

	if err := model.DB.Create(&comment).Error; err != nil {
		return serializer.CommentActionResponse{
			Response: serializer.Response{
				StatusCode: 1,
				StatusMsg:  "DataBase save comment failed",
			},
		}
	}

	return serializer.CommentActionResponse{
		Response: serializer.Response{StatusCode: 0},
		CommentResponse: model.CommentResponse{
			Id:         1,
			User:       user,
			Content:    content,
			CreateDate: time.Now().Format("01-02"),
		},
	}
}

func (service *CommentService) CommentCancle() serializer.CommentCancleResponse {
	return serializer.CommentCancleResponse{
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  "Cancle comment successfully",
		},
	}
}
