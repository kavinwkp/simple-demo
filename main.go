package main

import (
	"github.com/RaymondCode/simple-demo/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
