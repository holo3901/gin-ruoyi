package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ruoyi/pkg/live"
)

func GetPushUrl(ctx *gin.Context) {
	streamName := ctx.Query("streamName")
	if streamName == "" {
		zap.L().Error("ResumeDelayLiveStream", zap.Error(errors.New("直播流名称不能为空")))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errors.New("直播流名称不能为空"))
		return
	}
	l := live.NewLive()
	res := l.GetPushUrl(streamName)
	ResponseSuccess(ctx, res)
}

// 获取拉流地址
func GetPullUrl(ctx *gin.Context) {
	streamName := ctx.Query("streamName")
	if streamName == "" {
		zap.L().Error("ResumeDelayLiveStream", zap.Error(errors.New("直播流名称不能为空")))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errors.New("直播流名称不能为空"))
		return
	}
	l := live.NewLive()
	res := l.GetPullUrl(streamName)
	ResponseSuccess(ctx, res)

}
