package ctrl

import (
	"sms_server/common"
	"sms_server/r2m"

	"github.com/gin-gonic/gin"
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
