package controllers

import "github.com/gin-gonic/gin"

func Index(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"Code":    200,
		"Message": "Welcome to REST API APLIKASI EOQ",
	})
}
