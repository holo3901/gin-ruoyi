package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ruoyi/pkg/erweima"
	"strconv"
)

func GenerateArticlePoster(ctx *gin.Context) {
	address := ctx.Param("address")
	id := ctx.Param("tupianid")
	ids, _ := strconv.Atoi(id)
	bytes, err := erweima.Erweima(address, ids)
	if err != nil {
		zap.L().Error("erweima.Erweima failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ctx.Header("content=Type", "image/pnglcharset=utf-8")
	ctx.String(200, string(bytes))
}
