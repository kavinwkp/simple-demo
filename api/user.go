package api

import (
	"fmt"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var userRegister = service.UserService{
		UserName: username,
		Password: password,
	}
	c.JSON(200, userRegister.Register())
	return
}

func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var userLogin = service.UserService{
		UserName: username,
		Password: password,
	}
	fmt.Println(userLogin)
	c.JSON(200, userLogin.Login())
	return
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	var userInfo = service.UserInfoService{
		Token: token,
	}
	c.JSON(200, userInfo.Info())
	return
}
