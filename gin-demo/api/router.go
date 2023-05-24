package api

import (
	"gin-demo/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	r.Use(middleware.CORS())

	rGroup := r.Group("/login")
	{
		rGroup.POST("", login)

		rGroup.POST("/username", UserName)
	}

	r.POST("/register", register) // 注册

	r.Run(":8080")
}
