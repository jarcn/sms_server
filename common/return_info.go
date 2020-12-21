package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReturnInfo struct {
	Description string
	RetCode     int
	RetObj      interface{}
	Success     bool
}

func Failed(c *gin.Context) {
	returnInfo := ReturnInfo{
		Description: "failed",
		RetCode:     0,
		RetObj:      nil,
		Success:     false,
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": returnInfo,
	})
}

func Success(c *gin.Context) {
	returnInfo := ReturnInfo{
		Description: "success",
		RetCode:     0,
		RetObj:      nil,
		Success:     true,
	}
	c.JSON(http.StatusOK, gin.H{
		"ret": returnInfo,
	})
}

func Error(c *gin.Context, info ReturnInfo) {
	c.JSON(http.StatusOK, gin.H{
		"ret": info,
	})
}
