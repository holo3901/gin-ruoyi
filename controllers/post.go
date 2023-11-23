package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"ruoyi/logic"
	"ruoyi/models"
	"ruoyi/pkg/exce"
	"strconv"
)

func ListPost(ctx *gin.Context) {
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

	var searchDictData *models.SysPost
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
	list, err := logic.SelectSysPostList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.selectsyspostlist invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportPost(ctx *gin.Context) {
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

	var searchDictData *models.SysPost
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
	list, err := logic.SelectSysPostList(p, searchDictData)
	if err != nil {
		zap.L().Error("logic.selectsyspostlist invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "postId",
		"title":  "岗位序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "postCode",
		"title":  "岗位编码",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "postName",
		"title":  "岗位名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "postSort",
		"title":  "岗位排序",
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
			var statusStr = ""
			if v.Status == "0" {
				statusStr = "正常"
			} else {
				statusStr = "停用"
			}
			data = append(data, map[string]interface{}{
				"postId":   v.PostId,
				"postCode": v.PostCode,
				"postName": v.PostName,
				"postSort": v.PostSort,
				"status":   statusStr,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "导出成功")
}

func GetPostInfo(ctx *gin.Context) {
	postId := ctx.Param("postId")
	id, _ := strconv.Atoi(postId)
	byId, err := logic.FindPostInfoById(id)
	if err != nil {
		zap.L().Error("logic.FindPostInfoById failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, byId)
}

func SavePost(ctx *gin.Context) {
	p := new(models.SysPost)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("savepost invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	userid, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.SavePost(userid, p)
	if err != nil {
		zap.L().Error("logic.savepost failed:", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}

func UploadPost(ctx *gin.Context) {
	p := new(models.SysPost)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("savepost invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	userid, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	err = logic.UploadPost(userid, p)
	if err != nil {
		zap.L().Error("logic.uploadpost failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "更新成功")
}

func DeletePost(ctx *gin.Context) {
	postid := ctx.Param("postIds")
	postId, _ := strconv.Atoi(postid)
	err := logic.DeletePost(postId)
	if err != nil {
		zap.L().Error("logic.deletepost failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除成功")
}

func GetPostOptionSelect(ctx *gin.Context) {
	param := new(models.SearchTableDataParam)
	p := new(models.SysPost)
	list, err := logic.SelectSysPostList(param, p)
	if err != nil {
		zap.L().Error("logic.selectSyspostlist failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}
