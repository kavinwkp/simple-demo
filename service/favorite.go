package service

import (
	"errors"
	"fmt"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

type FavoriteActionService struct {
	Token   string
	VideoId int64
}

func (service *FavoriteActionService) FavoriteAction() serializer.FavoriteActionResponse {
	fmt.Println("favorite action")
	claim, _ := utils.ParseToken(service.Token)
	video_id := service.VideoId
	user_id := claim.Id

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
	fmt.Println(favorite)

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
