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

func ListDict(ctx *gin.Context) {
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

	var searchDictData *models.SearchDictData
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
	list, err := logic.SelectDictDataList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.selectDictdata with error", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, gin.H{
		"data": list,
		"code": "success",
	})
}

func ExportDict(ctx *gin.Context) {
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

	var searchDictData *models.SearchDictData
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

	list, _ := logic.SelectDictDataList(p, searchDictData)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "dictCode",
		"title":  "字典编码",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictSort",
		"title":  "字典排序",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictLabel",
		"title":  "字典标签",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictValue",
		"title":  "字典键值",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictType",
		"title":  "字典类型",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "isDefault",
		"title":  "是否默认",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "状态",
		"width":  "10",
		"is_num": "0",
	})
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			defaults := v.IsDefault
			var df = ""
			if "Y" == defaults {
				df = "是"
			}
			if "N" == defaults {
				df = "否"
			}
			var status = v.Status
			statusStr := ""
			if status == "0" {
				statusStr = "正常"
			}
			if status == "1" {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"dictCode":  v.DictCode,
				"dictSort":  v.DictSort,
				"dictLabel": v.DictLabel,
				"dictValue": v.DictValue,
				"dictType":  v.DictType,
				"isDefault": df,
				"status":    statusStr,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)

	ResponseSuccess(ctx, "导出成功")
}

func GetDictCode(ctx *gin.Context) {
	dictcode := ctx.Param("dictCode")
	result, err := logic.FindDictCodeByID(dictcode)
	if err != nil {
		zap.L().Error("logic.FindDictCodeByID failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, result)
}

func DictTypeHandler(ctx *gin.Context) {
	dictType := ctx.Param("dictType")
	byType, err := logic.FindDictCodeByType(dictType)
	if err != nil {
		zap.L().Error("logic.FindDictCodeByType failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, byType)
}

func SaveDictData(ctx *gin.Context) {
	dictDataParam := new(models.SysDictData)
	if err := ctx.ShouldBindJSON(&dictDataParam); err != nil {
		zap.L().Error("savedictdata with invalid", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {

			ResponseError(ctx, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		zap.L().Error("用户未登录")
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.SaveDictData(id, dictDataParam)
	if err != nil {
		zap.L().Error("logic.savedictdata invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "success")
}

func UpDictData(ctx *gin.Context) {
	dictDataParam := new(models.SysDictData)
	if err := ctx.ShouldBindJSON(&dictDataParam); err != nil {
		zap.L().Error("savedictdata with invalid", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {

			ResponseError(ctx, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		zap.L().Error("logic.savedictdata invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	err = logic.SaveDictData(id, dictDataParam)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "更新成功")
}

func DeleteDictData(ctx *gin.Context) {
	dictcode := ctx.Param("dictCodes")

	err := logic.DeleteDictData(dictcode)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}

func ListDictType(ctx *gin.Context) {
	p := new(models.SearchTableDataParam)

	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("ListDictType invalid param", zap.Error(err))
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

	var searchDictType *models.SysDictType
	if err := json.Unmarshal(otherData, &searchDictType); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	list, err := logic.SelectSysDictTypeList(p, searchDictType)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportType(ctx *gin.Context) {
	p := new(models.SearchTableDataParam)

	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("ListDictType invalid param", zap.Error(err))
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

	var searchDictType *models.SysDictType
	if err := json.Unmarshal(otherData, &searchDictType); err != nil {
		zap.L().Error("listDict with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	list, err := logic.SelectSysDictTypeList(p, searchDictType)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "dictId",
		"title":  "字典主键",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictName",
		"title":  "字典名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "dictType",
		"title":  "字典类型",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "状态",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			status := v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "正常"
			}
			if status == "1" {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"dictId":   v.DictId,
				"dictName": v.DictName,
				"dictType": v.DictType,
				"status":   statusStr,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "导出成功")
}

func GetTypeDict(ctx *gin.Context) {
	dictId := ctx.Param("dictId")
	id, err := logic.FindTypeDictById(dictId)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, id)
}

func SaveType(ctx *gin.Context) {
	p := new(models.SysDictType)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("saveType invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		zap.L().Error("用户未登录")
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.SaveType(p, id)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存type成功")
}

func UploadType(ctx *gin.Context) {
	p := new(models.SysDictType)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("saveType invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}

		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		zap.L().Error("用户未登录")
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.SaveType(p, id)
}

func DeleteDataType(ctx *gin.Context) {
	dictIds := ctx.Param("dictIds")

	err := logic.DeleteDataType(dictIds)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}

func RefreshCache(ctx *gin.Context) {
	err := redis.DeleteDict()
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "成功")
}

func GetOptionSelect(ctx *gin.Context) {
	optionSelect, err := logic.GetOptionSelect()
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, optionSelect)
}
