package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	// 传入token video_id action_type
	token := c.Query("token")
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64) // string to int64
	action_type := c.Query("action_type")
	var favoriteServer = service.FavoriteActionService{
		Token:   token,
		VideoId: video_id,
	}
	switch action_type {
	case "1":
		c.JSON(200, favoriteServer.FavoriteAction())
	case "2":
		c.JSON(200, favoriteServer.FavoriteCancleAction())
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, serializer.VideoListResponse{
		Response: serializer.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
