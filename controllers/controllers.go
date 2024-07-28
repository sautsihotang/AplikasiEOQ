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
			"message": "succes create barang",
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
			"message": "succes - Get all barang",
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
			"code":     400,
			"message":  err.Error(),
			"supplier": barang,
		})
	} else {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "succes - Get barang",
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
			"message": "succes - delete barang",
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
			"message":  "succes create supplier",
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
			"message":  "succes - Get all supplier",
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
			"message":  "succes - Get supplier",
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
			"message": "succes - delete supplier",
		})
	}
}
