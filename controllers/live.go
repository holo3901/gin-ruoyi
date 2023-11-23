package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ruoyi/models"
	"ruoyi/pkg/live"
	"strconv"
)

// 3.查询直播中的流程
func GetLiveStreamOnlineList(ctx *gin.Context) {
	pageNum, _ := strconv.ParseUint(ctx.DefaultQuery("page_num", "1"), 10, 64)
	pagesize, _ := strconv.ParseUint(ctx.DefaultQuery("page_size", "10"), 10, 64)
	streamName := ctx.Query("stream_name")
	l := live.NewLive()
	streamNameList := make([]string, 0)
	if streamName != "" {
		streamNameList = append(streamNameList, streamName)
	}
	res, err := l.GetLiveStreamOnlineList(pageNum, pagesize, streamNameList...)
	if err != nil {
		zap.L().Error("GetLiveStreamOnlineList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, gin.H{
		"msg": res,
	})
}

// 4.获取禁推流列表
func GetLiveForbidStreamList(ctx *gin.Context) {
	pageNum, _ := strconv.ParseUint(ctx.DefaultQuery("pageNum", "1"), 10, 64)
	pagesize, _ := strconv.ParseUint(ctx.DefaultQuery("pageSize", "10"), 10, 64)
	streamName := ctx.Query("streamName")
	l := live.NewLive()
	streamNameList := make([]string, 0)
	if streamName != "" {
		streamNameList = append(streamNameList, streamName)
	}
	res, err := l.GetLiveForbidStreamList(pageNum, pagesize, streamNameList...)
	if err != nil {
		zap.L().Error("GetLiveForbidStreamList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, res)
}

// 5.禁直播推流
func ForbidLiveStream(ctx *gin.Context) {
	req := new(models.ForbidLiveStreamReq)
	err := ctx.ShouldBind(&req)
	if err != nil {
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}

	l := live.NewLive()
	res, err := l.ForbidLiveStream(req.StreamName, req.ResumeTime, req.Reason)
	if err != nil {
		zap.L().Error("ForbidLiveStream", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)
}

// 6.断开流
func DropLiveStream(ctx *gin.Context) {
	streamName := ctx.Query("streamName")
	l := live.NewLive()
	res, err := l.DropLiveStream(streamName)
	if err != nil {
		zap.L().Error("DropLiveStream", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)
}

// 7.查询流状态
func GetLiveStreamState(ctx *gin.Context) {
	streamName := ctx.Query("streamName")
	l := live.NewLive()
	res, err := l.GetLiveStreamState(streamName)
	if err != nil {
		zap.L().Error("GetLiveStreamState", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)
}

// 8.恢复直播流
func ResumeLiveStream(ctx *gin.Context) {
	streamName := ctx.Query("streamName")
	l := live.NewLive()
	res, err := l.ResumeLiveStream(streamName)
	if err != nil {
		zap.L().Error("ResumeLiveStream", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)
}

// 9.获取延迟播放列表
func GetLiveDelayInfoList(ctx *gin.Context) {
	l := live.NewLive()
	res, err := l.GetLiveDelayInfoList()
	if err != nil {
		zap.L().Error("GetLiveDelayInfoList", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)
}

// 10.设置延迟直播
func AddDelayLiveStream(ctx *gin.Context) {
	req := new(models.AddDelayLiveStreamReq)
	err := ctx.ShouldBind(&req)
	if err != nil {
		zap.L().Error("GetLiveDelayInfoList", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	l := live.NewLive()
	res, err := l.AddDelayLiveStream(req.StreamName, req.DelayTime, req.ExpireTime)
	if err != nil {
		zap.L().Error("AddDelayLiveStream", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)

}

// 11.取消延迟直播
func ResumeDelayLiveStream(ctx *gin.Context) {
	streamName := ctx.Query("streamName")
	l := live.NewLive()
	res, err := l.ResumeDelayLiveStream(streamName)
	if err != nil {
		zap.L().Error("ResumeDelayLiveStream", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, res)
}
