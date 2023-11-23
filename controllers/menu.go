package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"io"
	"ruoyi/logic"
	"ruoyi/models"
	"strconv"
)

func ListMenu(ctx *gin.Context) {
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

	var searchDictData *models.SysMenu
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
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	userId, err := logic.SelectSysMenuByUserId(id, searchDictData, p)
	if err != nil {
		zap.L().Error("logic.SelectSysMenuByUserId invalidparam", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, userId)
}

func GetMenuInfo(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	atoi, err := strconv.Atoi(menuId)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	info, err := logic.GetMenuInfo(atoi)
	if err != nil {
		zap.L().Error("logic.GetMenuInfo invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, info)

}

func GetTreeSelect(ctx *gin.Context) {
	data, _ := io.ReadAll(ctx.Request.Body)
	io.NopCloser(bytes.NewBuffer(data))

	p := new(models.SysMenu)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("getTreeSelect invalid param error", zap.Error(err))
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
	treeSelect, err := logic.GetTreeSelect(id, p)
	if err != nil {
		zap.L().Error("logic.GetTreeSelect invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, treeSelect)
}

func TreeSelectByRole(ctx *gin.Context) {
	roleId := ctx.Param("roleId")
	roleid, _ := strconv.Atoi(roleId)
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	p := new(models.SysMenu)
	treeSelect, err := logic.GetTreeSelect(id, p)
	if err != nil {
		zap.L().Error("logic.GetTreeSelect invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	role, err := logic.GetMenuTreeSelectByRole(roleid)
	if err != nil {
		zap.L().Error("logic.GetMenuTreeSelectByRole invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	ResponseSuccess(ctx, gin.H{
		"menus":       treeSelect,
		"checkedkeys": role,
	})
}

func SaveMenu(ctx *gin.Context) {
	p := new(models.SysMenu)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("savemenu invalid param", zap.Error(err))
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
	err = logic.SaveMenu(id, p)
	if err != nil {
		zap.L().Error("logic.savemenu invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "保存成功")
}

func UploadMenu(ctx *gin.Context) {
	p := new(models.SysMenu)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("savemenu invalid param", zap.Error(err))
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
	err = logic.UpdateMenu(id, p)
	if err != nil {
		zap.L().Error("logic.UpdateMenu failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "修改成功")
}

func DeleteMenu(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	atoi, _ := strconv.Atoi(menuId)
	err := logic.DeleteMenu(atoi)
	if err != nil {
		zap.L().Error("logic.DeleteMenu faield", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除菜单失败")
}
