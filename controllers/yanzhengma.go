package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"ruoyi/pkg/yanzhengma"
)

func CaptchaImageHandler(ctx *gin.Context) {
	id, b64s, err := yanzhengma.CreateImageCaptcha()
	if err != nil {
		zap.L().Error("middleware.CreateImageCaptcha", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
	}
	ResponseSuccess(ctx, gin.H{
		"msg":            "操作成功",
		"img":            b64s,
		"code":           http.StatusOK,
		"captchaEnabled": true,
		"uuid":           id})
}
