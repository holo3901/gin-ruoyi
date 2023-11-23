package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/logic"
	"ruoyi/models"
	"ruoyi/pkg/exce"
	"strconv"
)

func ListConfig(ctx *gin.Context) {
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

	var searchDictData *models.SysConfig
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
	list, err := logic.SelectConfigList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.SelectConfigList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportConfig(ctx *gin.Context) {
	p := new(models.SysConfig)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("ExportConfig with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	param := new(models.SearchTableDataParam)

	list, err := logic.SelectConfigList(param, p)
	if err != nil {
		zap.L().Error("logic.SelectConfigList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "configId",
		"title":  "参数主键",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configName",
		"title":  "参数名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configKey",
		"title":  "参数键名",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configValue",
		"title":  "参数键值",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "configType",
		"title":  "系统内置",
		"width":  "10",
		"is_num": "0",
	})

	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var configTypeStr = ""
			if v.ConfigType == "Y" {
				configTypeStr = "是"
			}
			if v.ConfigType == "N" {
				configTypeStr = "否"
			}
			data = append(data, map[string]interface{}{
				"configId":    v.ConfigId,
				"configName":  v.ConfigName,
				"configKey":   v.ConfigKey,
				"configValue": v.ConfigValue,
				"configType":  configTypeStr,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "成功")
}

func GetConfigInfo(ctx *gin.Context) {
	configId := ctx.Param("configId")
	id, _ := strconv.Atoi(configId)
	info, err := logic.GetConfigInfo(id)
	if err != nil {
		zap.L().Error("logic.GetConfigInfo failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, info)

}

func GetConfigKey(ctx *gin.Context) {
	configkey := ctx.Param("config_key")
	a := new(models.SysConfig)
	a.ConfigKey = configkey
	fmt.Println(a)
	config, err := logic.SelectConfig(a)
	if err != nil {
		zap.L().Error("logic.GetConfigKey failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, config)
}

func SaveConfig(ctx *gin.Context) {
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	p := new(models.SysConfig)
	if err = ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SaveConfig with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	err = logic.SaveConfig(id, p)
	if err != nil {
		zap.L().Error("logic.SaveConfig failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}

func UploadConfig(ctx *gin.Context) {
	userid, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	a := new(models.SysConfig)
	if err = ctx.ShouldBindJSON(&a); err != nil {
		zap.L().Error("UploadConfig with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	err = logic.UploadConfig(userid, a)
	if err != nil {
		zap.L().Error("logic.UploadConfig failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "更新成功")
}

func DeletectConfig(ctx *gin.Context) {
	configid := ctx.Param("configIds")
	err := logic.DeleteConfig(configid)
	if err != nil {
		zap.L().Error("logic.DeleteConfig failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除config成功")
}

func DeleteCacheConfig(ctx *gin.Context) {
	GetCurrentUserID(ctx)
	/*
		加载缓存
		重复初始化
		func loadingConfigCache() {
			//var param = tools.SearchTableDataParam{}
			//SelectConfigList(param, false)
			/*重新赋值进去

		}

		func DelCacheConfig(refreshCache string) R.Result {
			/*删除所有缓存*/
	/*重复初始化
		loadingConfigCache()
		return R.ReturnSuccess("操作成功")
	}

	*/
}
