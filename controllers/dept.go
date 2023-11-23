package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/logic"
	"ruoyi/models"
	"strconv"
)

func ListDept(ctx *gin.Context) {
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

	var searchDictData *models.SysDept
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
	list, err := logic.GetDeptList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic,GetDeptList failed ", zap.Error(err))
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	ResponseSuccess(ctx, list)

}

func ExcludeDept(ctx *gin.Context) {
	deptId := ctx.Param("deptId")
	a := new(models.SearchTableDataParam)
	p := new(models.SysDept)
	dept, err := logic.ExcludeDept(deptId, a, p)
	if err != nil {
		zap.L().Error("logic.ExcludeDept failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, dept)
}

func GetDept(ctx *gin.Context) {
	deptId := ctx.Param("deptId")
	id, _ := strconv.Atoi(deptId)
	result, err := logic.GetDeptInfo(id)
	if err != nil {
		zap.L().Error("logic.GetDept failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, result)
}

func SaveDept(ctx *gin.Context) {
	p := new(models.SysDept)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("saveDept invalid.param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errs.Translate(trans))
		return
	}
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.SaveDept(p, id)
	if err != nil {
		zap.L().Error("logic.SaveDept failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}

func UpDataDept(ctx *gin.Context) {
	p := new(models.SysDept)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("saveDept invalid.param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errs.Translate(trans))
		return
	}
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.UploadDept(p, id)
	if err != nil {
		zap.L().Error("logic.UploadDept failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}

func DeleteDept(ctx *gin.Context) {
	deptId := ctx.Param("deptId")
	id, _ := strconv.Atoi(deptId)
	if err := logic.DeleteDept(id); err != nil {
		zap.L().Error("logic.DeleteDept failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除dept")
}
