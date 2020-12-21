package ctrl

import (
	"github.com/gin-gonic/gin"
	"sms_server/common"
	"sms_server/services"
)

type SendCodePara struct {
	PhoneNum string `form:"phoneNum" json:"phoneNum" binding:"required"`
	CodeType int    `form:"codeType" json:"codeType" binding:"required"`
}

type CheckCodePara struct {
	PhoneNum  string `form:"phoneNum" json:"phoneNum" binding:"required"`
	CheckCode string `form:"checkCode" json:"checkCode" binding:"required"`
}

func SendCode(ctx *gin.Context) {
	var para SendCodePara
	if err := ctx.BindJSON(&para); err != nil {
		common.Failed(ctx)
		return
	}
	if para.PhoneNum == "" {
		common.Failed(ctx)
		return
	}
	if _, err := services.SendCode(para.PhoneNum, para.CodeType); err != nil {
		info := common.ReturnInfo{
			Description: "短信验证码发送异常",
			RetCode:     500,
			RetObj:      err.Error(),
			Success:     false,
		}
		common.Error(ctx, info)
		return
	}
	common.Success(ctx)
	return
}

func CheckCode(ctx *gin.Context) {
	var para CheckCodePara
	if err := ctx.BindJSON(&para); err != nil {
		common.Failed(ctx)
		return
	}
	if para.PhoneNum == "" || para.CheckCode == "" {
		common.Failed(ctx)
		return
	}
	if okCode, err := services.CheckCode(para.PhoneNum, para.CheckCode); !okCode {
		info := common.ReturnInfo{
			Description: "验证码错误",
			RetCode:     400,
			RetObj:      err.Error(),
			Success:     false,
		}
		common.Error(ctx, info)
		return
	} else {
		common.Success(ctx)
		return
	}
}
