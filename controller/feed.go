package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var videos []Video
	DB.Model(&Video{}).Order("id desc").Limit(5).Find(&videos)
	for index := range videos {
		videos[index].PlayUrl = baseURL + videos[index].PlayUrl
		videos[index].CoverUrl = baseURL + videos[index].CoverUrl
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
