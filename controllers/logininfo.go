package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/dao/redis"
	"ruoyi/logic"
	"ruoyi/models"
	"ruoyi/pkg/exce"
)

func LoginInformListHandler(ctx *gin.Context) {
	p := new(models.SearchTableDataParam)

	// 解析SearchTableDataParam.Other字段为SearchDictData类型
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	var otherData json.RawMessage
	if err := json.Unmarshal(p.Other, &otherData); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	var searchDictData *models.SysLogininfor
	if err := json.Unmarshal(otherData, &searchDictData); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	list, err := logic.LoginInformList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.LoginInformlist failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportHandler(ctx *gin.Context) {
	p := new(models.SearchTableDataParam)

	// 解析SearchTableDataParam.Other字段为SearchDictData类型
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	var otherData json.RawMessage
	if err := json.Unmarshal(p.Other, &otherData); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	var searchDictData *models.SysLogininfor
	if err := json.Unmarshal(otherData, &searchDictData); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	list, err := logic.LoginInformList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.LoginInformlist failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "infoId",
		"title":  "序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "userName",
		"title":  "用户账号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "登录状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "ipaddr",
		"title":  "登录地址",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginLocation",
		"title":  "登录地点",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "browser",
		"title":  "浏览器",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "os",
		"title":  "操作系统",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "msg",
		"title":  "提示消息",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginTime",
		"title":  "访问时间",
		"width":  "50",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var status = v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "成功"
			} else {
				statusStr = "失败"
			}

			var loginTime = v.LoginTime.Format(models.TimeFormat)
			data = append(data, map[string]interface{}{
				"infoId":        v.InfoId,
				"userName":      v.UserName,
				"status":        statusStr,
				"ipaddr":        v.Ipaddr,
				"loginLocation": v.LoginLocation,
				"browser":       v.Browser,
				"os":            v.Os,
				"msg":           v.Msg,
				"loginTime":     loginTime,
			})
		}

	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "成功")
}

func DeleteByIdHandler(ctx *gin.Context) {
	infoids := ctx.Param("infoIds")
	err := logic.DeleteInfoId(infoids)
	if err != nil {
		zap.L().Error("logic.deleteInfoId failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}

func CleanHandler(ctx *gin.Context) {
	err := logic.ClearLoginLog()
	if err != nil {
		zap.L().Error("logic.clearLoginLog failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "清理成功")
}

func UnlockHandler(ctx *gin.Context) {
	username := ctx.Param("username")
	err := redis.UnlockByUserName(username)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "清缓成功")
}
