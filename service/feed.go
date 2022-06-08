package service

import (
	"time"

	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
)

type FeedService struct{}

func (service *FeedService) Feed() serializer.FeedResponse {
	var videosTable []model.VideoTable
	model.DB.Order("id desc").Limit(5).Find(&videosTable)

	var videos []model.Video
	for _, v := range videosTable {
		var user model.User
		model.DB.First(&user, v.UserID)
		var video = model.Video{
			Id:            v.Id,
			Title:         v.Title,
			Author:        user,
			PlayUrl:       config.BaseURL + v.PlayUrl,
			CoverUrl:      config.BaseURL + v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		}
		videos = append(videos, video)
	}

	return serializer.FeedResponse{
		Response:  serializer.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	}
}
