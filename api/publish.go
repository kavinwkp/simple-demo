package api

import (
	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(200, serializer.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	var publishServer = service.PublishService{
		Token:      token,
		Title:      title,
		FileHeader: data,
	}
	c.JSON(200, publishServer.Publish())
	return
}

func PublishList(c *gin.Context) {
	token := c.Query("token")
	var publishListServer = service.PublishListService{
		Token: token,
	}
	c.JSON(200, publishListServer.PublishList())
	return
}
