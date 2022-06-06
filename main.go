package main

import (
	"fmt"

	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	var user model.User
	result := model.DB.Model(&model.User{}).Where("name=?", "kavin").First(&user)
	fmt.Println(result.RowsAffected)
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
