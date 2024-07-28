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

		//supplier
		apiV1Group.POST("/supplier", controllers.CreateSupplier)
		apiV1Group.GET("/suppliers", controllers.GetSuppliers)
		apiV1Group.GET("/supplier", controllers.GetSupplierbyId)
		apiV1Group.PUT("/supplier", controllers.UpdateSupplier)
		apiV1Group.DELETE("/supplier", controllers.DeleteSupplier)
	}

	return r
}
