package controller

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	baseURL = "http://10.37.62.58:8080/static/"
)

type VideoListResponse struct {
	serializer.Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token") // token放在data里，所有用PostForm
	title := c.PostForm("title")
	claim, _ := utils.ParseToken(token)
	username := claim.UserName
	var user model.User

	if err := model.DB.Where("name=?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, serializer.UserLoginResponse{
				Response: serializer.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, serializer.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 视频保存成功就在Video表中插入视记录
	var video = model.Video{
		UserID:        user.Id,
		User:          user,
		Title:         title,
		PlayUrl:       finalName,
		CoverUrl:      "cover.png",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	}
	err = model.DB.Create(&video).Error
	if err != nil {
		c.JSON(http.StatusOK, serializer.UserLoginResponse{
			Response: serializer.Response{StatusCode: 1, StatusMsg: "Database save video failed"},
		})
	}
	c.JSON(http.StatusOK, serializer.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id") // user_id放在param里，所有用Query
	var videos []model.Video
	model.DB.Model(&model.Video{}).Preload("User").Where("user_id=?", user_id).Find(&videos)

	for index := range videos {
		videos[index].PlayUrl = baseURL + videos[index].PlayUrl
		videos[index].CoverUrl = baseURL + videos[index].CoverUrl
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: serializer.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
