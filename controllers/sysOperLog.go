package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/logic"
	"ruoyi/pkg/exce"

	"ruoyi/models"
	"strconv"
)

func ListOperlog(ctx *gin.Context) {
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

	var searchDictData *models.SysOperLog
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
	list, err := logic.GetOperLogList(p, searchDictData)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func DelectOperlog(ctx *gin.Context) {
	id := ctx.Param("operId")
	logid, _ := strconv.Atoi(id)
	err := logic.DeleteOperlog(logid)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}

func ClearOperlog(ctx *gin.Context) {
	err := logic.ClearOperlog()
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "成功")
}

func ExportOperlog(ctx *gin.Context) {
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

	var searchDictData *models.SysOperLog
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
	list, err := logic.GetOperLogList(p, searchDictData)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "operId",
		"title":  "操作序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "title",
		"title":  "操作模块",
		"width":  "15",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "businessType",
		"title":  "业务类型",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "method",
		"title":  "请求方法",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "requestMethod",
		"title":  "请求方式",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operatorType",
		"title":  "操作类别",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operName",
		"title":  "操作人员",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptName",
		"title":  "部门名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operUrl",
		"title":  "请求地址",
		"width":  "60",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operIp",
		"title":  "操作地址",
		"width":  "50",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operLocation",
		"title":  "操作地点",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operParam",
		"title":  "请求参数",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jsonResult",
		"title":  "返回参数",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "状态",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "errorMsg",
		"title":  "错误消息",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "operTime",
		"title":  "操作时间",
		"width":  "30",
		"is_num": "0",
	})

	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var businessType = v.BusinessType
			var businessTypeStr = ""
			if businessType == "0" {
				businessTypeStr = "其它"
			} else if businessType == "1" {
				businessTypeStr = "新增"
			} else if businessType == "2" {
				businessTypeStr = "修改"
			} else if businessType == "3" {
				businessTypeStr = "删除"
			} else if businessType == "4" {
				businessTypeStr = "授权"
			} else if businessType == "5" {
				businessTypeStr = "导出"
			} else if businessType == "6" {
				businessTypeStr = "导入"
			} else if businessType == "7" {
				businessTypeStr = "强退"
			} else if businessType == "8" {
				businessTypeStr = "生成代码"
			} else if businessType == "9" {
				businessTypeStr = "清空数据"
			}
			var operatorType = v.OperatorType
			var operatorTypeStr = ""
			if operatorType == "0" {
				operatorTypeStr = "其它"
			} else if operatorType == "1" {
				operatorTypeStr = "后台用户"
			} else if operatorType == "2" {
				operatorTypeStr = "手机端用户"
			}
			var status = v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "正常"
			} else {
				statusStr = "停用"
			}
			var operTime = v.OperTime.Format(models.TimeFormat)
			data = append(data, map[string]interface{}{
				"operId":        v.OperId,
				"title":         v.Title,
				"businessType":  businessTypeStr,
				"method":        v.Method,
				"requestMethod": v.RequestMethod,
				"operatorType":  operatorTypeStr,
				"operName":      v.OperName,
				"deptName":      v.DeptName,
				"operUrl":       v.OperUrl,
				"operIp":        v.OperIp,
				"operLocation":  v.OperLocation,
				"operParam":     v.OperParam,
				"jsonResult":    v.JsonResult,
				"status":        statusStr,
				"errorMsg":      v.ErrorMsg,
				"operTime":      operTime,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "创建成功")
}
