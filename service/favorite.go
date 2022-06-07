package service

import (
	"errors"

	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

type FavoriteActionService struct {
	Token   string
	VideoId int64
}

type FavoriteListService struct {
	Token  string
	UserId int64
}

func (service *FavoriteActionService) FavoriteAction() serializer.FavoriteActionResponse {

	claim, _ := utils.ParseToken(service.Token)
	video_id := service.VideoId
	user_id := claim.UserId

	var favorite model.Favorite

	result := model.DB.Model(&model.Favorite{}).Where("user_id=? AND video_id=?", user_id, video_id).First(&favorite)
	if result.RowsAffected == 1 {
		return serializer.FavoriteActionResponse{
			Response: serializer.Response{
				StatusCode: 1,
				StatusMsg:  "Already favorite",
			},
		}
	}
	// 没点过赞就新增一条记录
	favorite.UserID = user_id
	favorite.VideoId = video_id

	err := model.DB.Create(&favorite).Error
	if err != nil {
		return serializer.FavoriteActionResponse{
			Response: serializer.Response{
				StatusCode: 1,
				StatusMsg:  "DataBase save favorite failed",
			},
		}
	}
	// 标记视频已点赞和更新点赞数
	model.DB.Model(&model.Video{}).Where("id = ?", video_id).Update("is_favorite", true)
	model.DB.Model(&model.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count+1"))

	return serializer.FavoriteActionResponse{
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  "Favorite successfully",
		},
	}
}

func (service *FavoriteActionService) FavoriteCancleAction() serializer.FavoriteActionResponse {
	claim, _ := utils.ParseToken(service.Token)
	video_id := service.VideoId
	user_id := claim.Id

	var favorite model.Favorite

	if err := model.DB.Model(&model.Favorite{}).Where("user_id=? AND video_id=?", user_id, video_id).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.FavoriteActionResponse{
				Response: serializer.Response{StatusCode: 1, StatusMsg: "Not favorite"},
			}
		}
	}
	model.DB.Delete(&favorite)
	model.DB.Model(&model.Video{}).Where("id = ?", video_id).Update("is_favorite", false)
	model.DB.Model(&model.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count-1"))

	return serializer.FavoriteActionResponse{
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  "Cancle favorite successfully",
		},
	}
}

func (service *FavoriteListService) FavoriteList() serializer.VideoListResponse {
	var favorites []model.Favorite
	model.DB.Model(&model.Favorite{}).Where("user_id=?", service.UserId).Find(&favorites)
	var video_ids []int64
	for _, v := range favorites {
		video_ids = append(video_ids, v.VideoId)
	}
	var videos []model.Video
	model.DB.Find(&videos, video_ids)
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
