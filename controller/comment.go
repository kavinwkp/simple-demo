package controller

import (
	"net/http"

	"github.com/RaymondCode/simple-demo/serializer"
	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	serializer.Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	serializer.Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: serializer.Response{StatusCode: 0},
				Comment: Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, serializer.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, serializer.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    serializer.Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
