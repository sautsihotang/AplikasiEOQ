package controllers

import (
	"aplikasieoq/models"
	"aplikasieoq/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "Welcome to REST API APLIKASI EOQ",
	})
}

// hitung stock barang
func CalculateStock(ctx *gin.Context) {
	stocks, err := service.CalculateStock(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"stock":   stocks,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - Get all stock barang",
			"stock":   stocks,
		})
	}
}

// hitung eoq
func CalculateEOQ(ctx *gin.Context) {
	eoq, err := service.CalculateEOQ(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"eoq":     eoq,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success create eoq",
			"eoq":     eoq,
		})
	}
}

// controller get all eoq
func GetEOQ(ctx *gin.Context) {
	eoq, err := service.GetEOQ(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"eoq":     eoq,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - Get all eoq",
			"eoq":     eoq,
		})
	}
}

func DeleteEoq(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeleteEoq(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete Eoq",
		})
	}
}

// controller create penjualan
func CreatePenjualan(ctx *gin.Context) {
	penjualan, err := service.CreatePenjualan(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":      400,
			"message":   err.Error(),
			"penjualan": penjualan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":      200,
			"message":   "success create penjualan",
			"penjualan": penjualan,
		})
	}
}

// controller get all penjualan
func GetPenjualans(ctx *gin.Context) {
	penjualan, err := service.GetPenjualans(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":      400,
			"message":   err.Error(),
			"penjualan": penjualan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":      200,
			"message":   "success - Get all penjualan",
			"penjualan": penjualan,
		})
	}
}

// controller penjualan by id
func GetPenjualanbyId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	penjualan, err := service.GetPenjualanbyId(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":      400,
			"message":   err.Error(),
			"penjualan": penjualan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":      200,
			"message":   "success - Get penjualan",
			"penjualan": penjualan,
		})
	}
}

func UpdatePenjualan(ctx *gin.Context) {
	var updatePenjualan models.Penjualan
	if err := ctx.ShouldBindJSON(&updatePenjualan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	penjualan, err := service.UpdatePenjualan(updatePenjualan)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "Penjualan not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"message":   "penjualan updated successfully",
		"penjualan": penjualan,
	})
}

func DeletePenjualan(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeletePenjualan(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete penjualan",
		})
	}
}

// controller create pemesanan
func CreatePemesanan(ctx *gin.Context) {
	pemesanan, err := service.CreatePemesanan(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":      400,
			"message":   err.Error(),
			"pemesanan": pemesanan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":      200,
			"message":   "success create pemesanan",
			"pemesanan": pemesanan,
		})
	}
}

// controller get all pemesanan
func GetPemesanans(ctx *gin.Context) {
	pemesanan, err := service.GetPemesanans(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":      400,
			"message":   err.Error(),
			"pemesanan": pemesanan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":      200,
			"message":   "success - Get all pemesanan",
			"pemesanan": pemesanan,
		})
	}
}

// controller pemesanan by id
func GetPemesananbyId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	pemesanan, err := service.GetPemesananbyId(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":      400,
			"message":   err.Error(),
			"pemesanan": pemesanan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":      200,
			"message":   "success - Get pemesanan",
			"pemesanan": pemesanan,
		})
	}
}

func UpdatePemesanan(ctx *gin.Context) {
	var updatePemesanan models.Pemesanan
	if err := ctx.ShouldBindJSON(&updatePemesanan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	pemesanan, err := service.UpdatePemesanan(updatePemesanan)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "Pemesanan not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"message":   "pemesanan updated successfully",
		"pemesanan": pemesanan,
	})
}

func DeletePemesanan(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeletePemesanan(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete pemesanan",
		})
	}
}

func Login(ctx *gin.Context) {
	data, err := service.Login(ctx)

	// Check if the user ID is zero or null
	if data.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Username atau password Salah !",
			"data":    nil,
		})
		return
	}

	// Handle other potential errors
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success login",
		"data":    data,
	})
}

// controller create penyimpanan
func CreateUser(ctx *gin.Context) {
	user, err := service.CreateUser(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"user":    user,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success create user",
			"user":    user,
		})
	}
}

// controller get all barang
func GetUsers(ctx *gin.Context) {
	user, err := service.GetUsers(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"user":    user,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - Get all user",
			"user":    user,
		})
	}
}

// controller user by id
func GetUserbyId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	user, err := service.GetUserbyId(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"user":    user,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - Get user",
			"user":    user,
		})
	}
}

func UpdateUser(ctx *gin.Context) {
	var updateUser models.User
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	user, err := service.UpdateUser(updateUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user updated successfully",
		"user":    user,
	})
}

func DeleteUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeleteUser(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete user",
		})
	}
}

// controller create penyimpanan
func CreatePenyimpanan(ctx *gin.Context) {
	penyimpanan, err := service.CreatePenyimpanan(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":        400,
			"message":     err.Error(),
			"penyimpanan": penyimpanan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":        200,
			"message":     "success create penyimpanan",
			"penyimpanan": penyimpanan,
		})
	}
}

// controller get all barang
func GetPenyimpanans(ctx *gin.Context) {
	penyimpanan, err := service.GetPenyimpanans(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":        400,
			"message":     err.Error(),
			"penyimpanan": penyimpanan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":        200,
			"message":     "success - Get all penyimpanan",
			"penyimpanan": penyimpanan,
		})
	}
}

// controller penyimpanan by id
func GetPenyimpananbyId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	penyimpanan, err := service.GetPenyimpananbyId(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":        400,
			"message":     err.Error(),
			"penyimpanan": penyimpanan,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":        200,
			"message":     "success - Get penyimpanan",
			"penyimpanan": penyimpanan,
		})
	}
}

func UpdatePenyimpanan(ctx *gin.Context) {
	var updatePenyimpanan models.Penyimpanan
	if err := ctx.ShouldBindJSON(&updatePenyimpanan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	penyimpanan, err := service.UpdatePenyimpanan(updatePenyimpanan)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "Penyimpanan not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"message":     "penyimpanan updated successfully",
		"penyimpanan": penyimpanan,
	})
}

func DeletePenyimpanan(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeletePenyimpanan(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete penyimpanan",
		})
	}
}

// controller create barang
func CreateBarang(ctx *gin.Context) {
	barang, err := service.CreateBarang(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"barang":  barang,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success create barang",
			"barang":  barang,
		})
	}
}

// controller get all barang
func GetBarangs(ctx *gin.Context) {
	barangs, err := service.GetBarangs(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"barang":  barangs,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - Get all barang",
			"barang":  barangs,
		})
	}
}

// controller barang by id
func GetBarangbyId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	barang, err := service.GetBarangByID(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
			"barang":  barang,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - Get barang",
			"barang":  barang,
		})
	}
}

func UpdateBarang(ctx *gin.Context) {
	var updatedBarang models.Barang
	if err := ctx.ShouldBindJSON(&updatedBarang); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	barang, err := service.UpdateBarang(updatedBarang)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "Barang not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "barang updated successfully",
		"barang":  barang,
	})
}

func DeleteBarang(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeleteBarang(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete barang",
		})
	}
}

// controller supplier
func CreateSupplier(ctx *gin.Context) {
	supplier, err := service.CreateSupplier(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":     400,
			"message":  err.Error(),
			"supplier": supplier,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":     200,
			"message":  "success create supplier",
			"supplier": supplier,
		})
	}
}

// controller supplier
func GetSuppliers(ctx *gin.Context) {
	suppliers, err := service.GetSuppliers(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":     400,
			"message":  err.Error(),
			"supplier": suppliers,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":     200,
			"message":  "success - Get all supplier",
			"supplier": suppliers,
		})
	}
}

// controller supplier by id
func GetSupplierbyId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	supplier, err := service.GetSupplierByID(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":     400,
			"message":  err.Error(),
			"supplier": supplier,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":     200,
			"message":  "success - Get supplier",
			"supplier": supplier,
		})
	}
}

func UpdateSupplier(ctx *gin.Context) {
	var updatedSupplier models.Supplier
	if err := ctx.ShouldBindJSON(&updatedSupplier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	supplier, err := service.UpdateSupplier(updatedSupplier)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "Supplier not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "Supplier updated successfully",
		"supplier": supplier,
	})
}

func DeleteSupplier(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": "Invalid ID",
		})
		return
	}
	err = service.DeleteSupplier(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"code":    400,
			"message": err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success - delete supplier",
		})
	}
}
