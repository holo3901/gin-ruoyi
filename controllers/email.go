package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/logic"
	"ruoyi/models"
	"strconv"
)

func SendEmail(ctx *gin.Context) {
	p := new(models.SendEmail)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("sendemail failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errs.Translate(trans))
		return
	}
	id, _ := GetCurrentUserID(ctx)
	err := logic.SendEmail(p, id)
	if err != nil {
		zap.L().Error("logic.SendEmail failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "发送成功")
}

func YanzhengEmail(ctx *gin.Context) {
	p := ctx.Param("id")
	s, _ := strconv.Atoi(p)
	id, _ := GetCurrentUserID(ctx)
	err := logic.YanzhengEmail(s, id)
	if err != nil {
		zap.L().Error("logic.YanzhengEmail failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, "验证成功")
}
