package service

import (
	"time"

	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
)

type FeedService struct{}

func (service *FeedService) Feed() serializer.FeedResponse {
	var videos []model.Video
	model.DB.Model(&model.Video{}).Order("id desc").Limit(5).Find(&videos)
	for index := range videos {
		videos[index].PlayUrl = config.BaseURL + videos[index].PlayUrl
		videos[index].CoverUrl = config.BaseURL + videos[index].CoverUrl
	}
	return serializer.FeedResponse{
		Response:  serializer.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	}
}
