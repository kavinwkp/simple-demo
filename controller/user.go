package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
	// UserId   int64  `json:"user_id,omitempty"`
	// UserName string `json:"user_name"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var user User
	result := DB.Model(&User{}).Where("name=?", username).First(&user)
	if result.RowsAffected == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		// atomic.AddInt64(&userIdSequence, 1)
		user.Name = username
		user.SetPassword(password)
		err := DB.Create(&user).Error
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "DataBase save user failed"},
			})
		}
		token, _ := GenerateToken(user.Id, user.Name)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var user User

	if err := DB.Where("name=?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}
	} else { // 用户名没错就验证密码
		if user.CheckPassword(password) == false {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 1, StatusMsg: "Password error"},
			})
		} else { // 密码正确，签发token
			token, _ := GenerateToken(user.Id, user.Name)
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   user.Id,
				Token:    token,
			})
		}
	}

}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	claim, _ := ParseToken(token)
	username := claim.UserName
	var user User
	if err := DB.Where("name=?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
			})
		}
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     user,
	})
}
