package router

import (
	"aplikasieoq/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	// Konfigurasi CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	apiV1Group := r.Group("/api/v1")
	{
		//TEST
		apiV1Group.GET("/", controllers.Index)

		//EOQ
		apiV1Group.POST("/eoq", controllers.CalculateEOQ)
		apiV1Group.GET("/eoq/all", controllers.GetEOQ)
		apiV1Group.DELETE("/eoq", controllers.DeleteEoq)

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

		//user
		apiV1Group.POST("/user", controllers.CreateUser)
		apiV1Group.GET("/user/all", controllers.GetUsers)
		apiV1Group.GET("/user", controllers.GetUserbyId)
		apiV1Group.PUT("/user", controllers.UpdateUser)
		apiV1Group.DELETE("/user", controllers.DeleteUser)

		//pemesanan
		apiV1Group.POST("/pemesanan", controllers.CreatePemesanan)
		apiV1Group.GET("/pemesanan/all", controllers.GetPemesanans)
		apiV1Group.GET("/pemesanan", controllers.GetPemesananbyId)
		apiV1Group.PUT("/pemesanan", controllers.UpdatePemesanan)
		apiV1Group.DELETE("/pemesanan", controllers.DeletePemesanan)

		//penjualan
		apiV1Group.POST("/penjualan", controllers.CreatePenjualan)
		apiV1Group.GET("/penjualan/all", controllers.GetPenjualans)
		apiV1Group.GET("/penjualan", controllers.GetPenjualanbyId)
		apiV1Group.PUT("/penjualan", controllers.UpdatePenjualan)
		apiV1Group.DELETE("/penjualan", controllers.DeletePenjualan)
	}

	return r
}
