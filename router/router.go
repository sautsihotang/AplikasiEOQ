package router

import (
	"aplikasieoq/controllers"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	apiV1Group := r.Group("/api/v1")
	{
		//TEST
		apiV1Group.GET("/", controllers.Index)
	}

	return r
}
