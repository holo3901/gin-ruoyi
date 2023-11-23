package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/logic"
	"ruoyi/models"
	"ruoyi/pkg/exce"
	"ruoyi/scheduler"
	"strconv"
)

func ListJob(ctx *gin.Context) {
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

	var searchDictData *models.SysJob
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
	list, err := logic.SelectJobList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.SelectJobList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportJob(ctx *gin.Context) {
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

	var searchDictData *models.SysJob
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
	list, err := logic.SelectJobList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.SelectJobList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "jobId",
		"title":  "任务序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobName",
		"title":  "任务名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobGroup",
		"title":  "任务组名",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "invokeTarget",
		"title":  "调用目标字符串",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "cronExpression",
		"title":  "执行表达式",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "misfirePolicy",
		"title":  "计划策略",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "concurrent",
		"title":  "并发执行",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "任务状态",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			misfirePolicyKey := v.MisfirePolicy
			var misfirePolicy = ""
			if 0 == misfirePolicyKey {
				misfirePolicy = "默认"
			}
			if 1 == misfirePolicyKey {
				misfirePolicy = "立即触发执行"
			}
			if 2 == misfirePolicyKey {
				misfirePolicy = "触发一次执行"
			}
			if 3 == misfirePolicyKey {
				misfirePolicy = "不触发立即执行"
			}
			concurrentKey := v.Concurrent
			var concurrent = ""
			if 0 == concurrentKey {
				concurrent = "允许"
			}
			if 1 == concurrentKey {
				concurrent = "禁止"
			}
			statusKey := v.Concurrent
			var status = ""
			if 0 == statusKey {
				status = "正常"
			}
			if 1 == statusKey {
				status = "暂停"
			}
			data = append(data, map[string]interface{}{
				"jobId":          v.JobId,
				"jobName":        v.JobName,
				"jobGroup":       v.JobGroup,
				"cronExpression": v.CronExpression,
				"invokeTarget":   v.InvokeTarget,
				"misfirePolicy":  misfirePolicy,
				"concurrent":     concurrent,
				"status":         status,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "成功")
}

func GetJobById(ctx *gin.Context) {
	jobId := ctx.Param("jobId")
	id, _ := strconv.Atoi(jobId)
	byId, err := logic.FindJobById(id)
	if err != nil {
		zap.L().Error("logic.FindJobById Failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, byId)
}

func SaveJob(ctx *gin.Context) {
	p := new(models.SysJobParam)
	if err := ctx.ShouldBindJSON(&p); err != nil {
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
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.SaveJob(p, id)
	if err != nil {
		zap.L().Error("logic.SaveJob failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}
func UploadJob(ctx *gin.Context) {
	p := new(models.SysJobParam)
	if err := ctx.ShouldBindJSON(&p); err != nil {
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
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.UploadJob(p, id)
	if err != nil {
		zap.L().Error("logic.UploadJob failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}

func ChangeStatus(ctx *gin.Context) {
	jobId := ctx.Param("jobIds")
	status := ctx.Param("status")
	err := logic.ChangeStatus(jobId, status)
	if err != nil {
		zap.L().Error("logic.ChangeStatus failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "修改成功")
}

func RunJob(ctx *gin.Context) {
	jobId := ctx.Param("jobId")
	id, _ := strconv.Atoi(jobId)
	result, err := logic.FindJobById(id)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	scheduler.RunOne(result)
	ResponseSuccess(ctx, result.CronExpression)
}

func DelectJob(ctx *gin.Context) {
	jobIds := ctx.Param("jobIds")
	err := logic.DeleteJob(jobIds)
	if err != nil {
		zap.L().Error("logic.deletejob failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除工作")
}

func ListJobLog(ctx *gin.Context) {
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

	var searchDictData *models.SysJobLog
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
	list, err := logic.SelectJobLogList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.SelectJobLogList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportJobLog(ctx *gin.Context) {
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

	var searchDictData *models.SysJobLog
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
	list, err := logic.SelectJobLogList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.SelectJobLogList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "jobLogId",
		"title":  "日志序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobName",
		"title":  "任务名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobGroup",
		"title":  "任务组名",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "invokeTarget",
		"title":  "调用目标字符串",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobMessage",
		"title":  "日志信息",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "执行状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "exceptionInfo",
		"title":  "异常信息",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			statusStr := ""
			status := v.Status
			if status == "0" {
				statusStr = "正常"
			}
			if status == "1" {
				statusStr = "失败"
			}
			data = append(data, map[string]interface{}{
				"jobLogId":      v.JobLogId,
				"jobName":       v.JobName,
				"jobGroup":      v.JobGroup,
				"invokeTarget":  v.InvokeTarget,
				"jobMessage":    v.JobMessage,
				"status":        statusStr,
				"exceptionInfo": v.ExceptionInfo,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)

	ResponseSuccess(ctx, "成功")
}

func GetJobLog(ctx *gin.Context) {
	joblog := ctx.Param("jobLogIds")
	joblogid, _ := strconv.Atoi(joblog)
	id, err := logic.FindJobLogById(joblogid)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, id)
}

func DetectJobLog(ctx *gin.Context) {
	joblog := ctx.Param("jobLogIds")
	joblogid, _ := strconv.Atoi(joblog)
	err := logic.DeleteJobLog(joblogid)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}

func ClearJobLog(ctx *gin.Context) {
	err := logic.ClearJobLog()
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}
