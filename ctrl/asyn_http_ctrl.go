package ctrl

import (
	"github.com/gin-gonic/gin"
	"sms_server/common"
	"sms_server/r2m"
)

type Asynchronous struct {
	Key  string `form:"key" json:"key" binding:"required"`
	Data string `form:"data" json:"data" binding:"required"`
}

func AsynchronousPost(ctx *gin.Context) {
	var para Asynchronous
	if err := ctx.BindJSON(&para); err != nil {
		common.Failed(ctx)
		return
	}
	job := r2m.Job{
		Key:  para.Key,
		Data: para.Data,
	}
	r2m.JobChannel <- job
	common.Success(ctx)
	return
}
