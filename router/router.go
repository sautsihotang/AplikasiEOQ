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

		//barang
		apiV1Group.POST("/barang", controllers.CreateBarang)
		apiV1Group.GET("/barangs", controllers.GetBarangs)
		apiV1Group.GET("/barang", controllers.GetBarangbyId)
		apiV1Group.PUT("/barang", controllers.UpdateBarang)
		apiV1Group.DELETE("/barang", controllers.DeleteBarang)

		//penyimpanan
		apiV1Group.POST("/penyimpanan", controllers.CreatePenyimpanan)
		apiV1Group.GET("/penyimpanan/all", controllers.GetPenyimpanans)
		apiV1Group.GET("/penyimpanan", controllers.GetPenyimpananbyId)
		apiV1Group.PUT("/penyimpanan", controllers.UpdatePenyimpanan)
		apiV1Group.DELETE("/penyimpanan", controllers.DeletePenyimpanan)
	}

	return r
}
