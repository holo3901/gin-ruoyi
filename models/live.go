package models

import "time"

type ForbidLiveStreamReq struct {
	StreamName string    `form:"stream_name" json:"streamName" binding:"required"`
	ResumeTime time.Time `form:"resume_time" json:"resumeTime" binding:"required" time_format:"2006-01-02 15:04:05"`
	Reason     string    `form:"reason" json:"reason" binding:"required"`
}

type AddDelayLiveStreamReq struct {
	StreamName string    `form:"stream_name" json:"streamName" binding:"required"`
	ExpireTime time.Time `form:"expire_time" json:"expireTime" binding:"required" time_format:"2006-01-02 15:04:05"`
	DelayTime  uint64    `form:"delay_time" json:"delayTime" binding:"required"`
}
