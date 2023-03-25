package main

import (
	"github.com/BeanWei/BWOnlineMusicPlayer/common"
	"github.com/BeanWei/BWOnlineMusicPlayer/config"
	"github.com/BeanWei/BWOnlineMusicPlayer/controllers"
	"github.com/gin-gonic/gin"
)


func main() {
	router := gin.New()
	// Gin原生中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 自定义的中间件
	router.Use(common.CORSMiddleware())

	// 静态文件
	router.LoadHTMLFiles("templates/index.html")
	router.Static("/static", "./static")
	router.Static("/images", "./static/images")
	router.Static("/music", config.MediaPath)

	// 请求
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// api group v1
	v1 := router.Group("/api/v1")
	{
		v1.POST("/music", controllers.MusicApiHandler)
	}

	router.Run(":3000")
}
