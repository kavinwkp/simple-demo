package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

type PublishService struct {
	Token      string
	Title      string
	FileHeader *multipart.FileHeader
}

type PublishListService struct {
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}

func (service *PublishService) Publish() serializer.PublishResponse {
	claim, _ := utils.ParseToken(service.Token)

	var user model.User
	username := claim.UserName
	if err := model.DB.Where("name=?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.PublishResponse{
				Response: serializer.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			}
		}
	}
	filename := filepath.Base(service.FileHeader.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := utils.UploadVideo(service.FileHeader, saveFile); err != nil {
		return serializer.PublishResponse{
			Response: serializer.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		}
	}
	// 视频保存成功就在Video表中插入视记录
	var video = model.Video{
		UserID:        user.Id,
		Title:         service.Title,
		PlayUrl:       finalName,
		CoverUrl:      "cover.png",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	}
	err := model.DB.Create(&video).Error
	if err != nil {
		return serializer.PublishResponse{
			Response: serializer.Response{StatusCode: 1, StatusMsg: "Database save video failed"},
		}
	}
	return serializer.PublishResponse{
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		},
	}
}

func (service *PublishListService) PublishList() serializer.VideoListResponse {
	var videos []model.Video
	model.DB.Model(&model.Video{}).Where("user_id=?", service.UserID).Find(&videos)

	for index := range videos {
		videos[index].PlayUrl = config.BaseURL + videos[index].PlayUrl
		videos[index].CoverUrl = config.BaseURL + videos[index].CoverUrl
	}

	return serializer.VideoListResponse{
		Response: serializer.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	}
}
