package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ruoyi/dao/redis"
	"ruoyi/pkg/exce"
	"ruoyi/pkg/yanzhengma"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"

	"ruoyi/dao/mysql"
	"ruoyi/logic"
	"ruoyi/models"
)

func LoginHandler(ctx *gin.Context) {
	p := new(models.LoginParam)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	if !yanzhengma.VerifyCaptcha(p.Uuid, p.Code) {
		zap.L().Error("验证码错误", zap.Error(errors.New("验证码错误")))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, "验证码错误")
		return
	}
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.signup failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
			return
		}

		ResponseError(ctx, CodeInvalidPassword)
		return
	}

	mysql.LoginInfoAdd(ctx, p, user, "登录成功", true)

	id := GetOnlineUserCount()
	ResponseSuccess(ctx, gin.H{
		"user_id":   fmt.Sprintf(":%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
		"onlineNum": id,
	})
}

func GetInfoHandler(ctx *gin.Context) {
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	info, err := logic.GetUserInfo(id)
	var a []string
	for i := 0; i < len(info.Rolepremission); i++ {
		a = append(a, info.Rolepremission[i].RoleName)
	}
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, gin.H{
		"user": gin.H{
			"username":     info.UserName,
			"nickname":     info.NickName,
			"phonenumbeer": info.Phonenumber,
			"email":        info.SysUser.Email,
			"avatar":       info.Avatar,
			"sex":          info.Sex,
			"createTime":   info.SysUser.CreateTime,
			"dept":         info.SysDept,
		},
		"roles":       a,
		"permissions": info.MenuPremission,
	})
}

func LogoutHandler(ctx *gin.Context) {
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	user, _ := mysql.FindUserById(id)
	err = mysql.DeleteOnlineByName(user.UserName)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	redis.CreateLogin(id, "") //将token变为0
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusFound, "/login")
}

func GetRoutersHandler(ctx *gin.Context) {
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	routers, err := logic.GetRouters(id)
	if err != nil {
		zap.L().Error("logic.GetRouters", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, routers)
}

func UpdatePwdHandler(ctx *gin.Context) {
	p := new(models.Repassword)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("login with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	id, _ := GetCurrentUserID(ctx)
	err := logic.UpdatePwd(id, p.OldPassword, p.NewPassword)
	if err != nil {
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, "修改成功")
}

// 查询个人信息
func ProfileHandler(ctx *gin.Context) {
	id, _ := GetCurrentUserID(ctx)
	info, err := logic.GetUserInfo(id)
	var a []string
	for i := 0; i < len(info.Rolepremission); i++ {
		a = append(a, info.Rolepremission[i].RoleName)
	}
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	ResponseSuccess(ctx, gin.H{
		"data": gin.H{
			"username":    info.UserName,
			"nickname":    info.NickName,
			"phonenumber": info.Phonenumber,
			"email":       info.SysUser.Email,
			"sex":         info.Sex,
			"createTime":  info.SysUser.CreateTime.Format(models.TimeFormat),
		},
		"roleGroup": a,
		"potGroup":  info.PostGroup,
	})
}

func PostProfileHandler(ctx *gin.Context) {
	p := new(models.Userparam)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("PostProfile with invalid param", zap.Error(err))
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
		return
	}
	err = logic.EditProfileUserInfo(id, p)
	if err != nil {
		zap.L().Error("logic.EditProfileUserInfo", zap.Error(err))
		ResponseError(ctx, CodeUserNotExist)
		return
	}
	ResponseSuccess(ctx, "修改成功")
}

func AvatarHandler(ctx *gin.Context) {
	file, fileHeader, err := ctx.Request.FormFile("file")
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	fileSize := fileHeader.Size
	id, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	err = logic.UpdateAvatar(file, fileSize, fileHeader.Filename, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeServerBusy, err.Error())
		return
	}
	ResponseSuccess(ctx, gin.H{
		"info": "头像更新成功",
	})
}

// ----------------------------- 用户管理 ------------------------------

func ListUser(ctx *gin.Context) {
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

	var searchDictType *models.SysUser
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
	list, err := logic.SelectUserList(p, searchDictType)
	if err != nil {
		ResponseErrorWithMsg(ctx, CodeInvalidParam, err)
		return
	}
	ResponseSuccess(ctx, list)
}

func ExportExport(ctx *gin.Context) {
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

	var searchDictType *models.SysUser
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

	list, err := logic.SelectUserParmList(p, searchDictType)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "userId",
		"title":  "用户序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptId",
		"title":  "部门编号",
		"width":  "15",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "userName",
		"title":  "登录名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "nickName",
		"title":  "用户名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "email",
		"title":  "用户邮箱",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "phonenumber",
		"title":  "手机号码",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "sex",
		"title":  "用户性别",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "帐号状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginIp",
		"title":  "最后登录IP",
		"width":  "30",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "loginDate",
		"title":  "最后登录时间",
		"width":  "60",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptName",
		"title":  "部门名称",
		"width":  "50",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "leader",
		"title":  "部门负责人",
		"width":  "30",
		"is_num": "0",
	})

	//填充数据
	data := make([]map[string]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			var sexStatus = v.Sex
			var sex = ""
			if sexStatus == "0" {
				sex = "男"
			} else if sexStatus == "1" {
				sex = "女"
			} else {
				sex = "未知"
			}
			var status = v.Status
			var statusStr = ""
			if status == "0" {
				statusStr = "正常"
			} else {
				statusStr = "停用"
			}
			var loginData = v.LoginDate.Format(models.TimeFormat)
			data = append(data, map[string]interface{}{
				"userId":      v.UserId,
				"deptId":      v.DeptId,
				"username":    v.UserName,
				"nickname":    v.NickName,
				"email":       v.Email,
				"phonenumber": v.Phonenumber,
				"sex":         sex,
				"status":      statusStr,
				"loginIp":     v.LoginIp,
				"loginDate":   loginData,
				"deptName":    v.DeptName,
				"leader":      v.Leader,
			})
		}
	}
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "导出成功")
}

func ImportUserData(ctx *gin.Context) {
	file, _, errload := ctx.Request.FormFile("file")
	if errload != nil {
		zap.L().Error("获取文件上传错误：", zap.Error(errload))
		ResponseErrorWithMsg(ctx, CodeInvalidParam, "获取上传文件错误"+errload.Error())
		return
	}
	xlsx, err := excelize.OpenReader(file) //打开读取到的文件
	if err != nil {
		zap.L().Error("ImportUserData，请选择文件", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeServerBusy, "请选择文件")
		return
	}

	updateSupport := ctx.Param("updateSupport")
	rows, _ := xlsx.GetRows("sheet1")
	users := make([]*models.SysUserParam, 0, 2)
	for irow, row := range rows {
		if irow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}
			atoi, _ := strconv.Atoi(data[1])
			id, _ := strconv.Atoi(data[0])
			users = append(users, &models.SysUserParam{
				UserId:      id,
				UserName:    data[2],
				NickName:    data[3],
				Email:       data[4],
				Phonenumber: data[5],
				Sex:         data[6],
				Status:      data[7],
				CreateTime:  time.Now(),
				DeptId:      atoi,
				PostIds:     logic.Split(data[8]),
				RoleIds:     logic.Split(data[9]),
			})
		}
	}
	if len(users) == 0 {
		zap.L().Error("ImportUserData,没有数据")
		ResponseErrorWithMsg(ctx, CodeInvalidParam, "请在表格中添加数据")
		return
	}
	id, _ := GetCurrentUserID(ctx)
	data, errs := logic.ImportUserData(users, updateSupport, id)
	if errs != nil {
		ResponseErrorWithMsg(ctx, CodeInvalidParam, errs)
		return
	}
	ResponseSuccess(ctx, data)
}

// 下载模板

func ImportTemplate(ctx *gin.Context) {
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "userId",
		"title":  "用户序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "deptId",
		"title":  "部门编号",
		"width":  "15",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "userName",
		"title":  "登录名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "nickName",
		"title":  "用户名称",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "email",
		"title":  "用户邮箱",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "phonenumber",
		"title":  "手机号码",
		"width":  "20",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "sex",
		"title":  "用户性别",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "帐号状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "岗位",
		"width":  "11",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "权限",
		"width":  "12",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	ex := exce.NewMyExcel()
	ex.ExportToWeb(dataKey, data, ctx)
	ResponseSuccess(ctx, "导出成功,也不知道导哪去了")
}

func GetUserInfo(ctx *gin.Context) {
	//参数用户
	useridStr := ctx.Param("userId")
	userid, err := GetCurrentUserID(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	if useridStr == "" {
		useridStr = strconv.Itoa(int(userid))
	}
	useridP, _ := strconv.Atoi(useridStr)

	err = logic.CheckUserDataScope(userid, useridStr)
	if err != nil {
		zap.L().Error("logic.CheckUserDataScope failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	info, err := logic.GetUserInfo(int64(useridP))
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	roles, err := logic.SelectRolePermissionByUserId(userid)
	if err != nil {
		zap.L().Error("logic.SelectRolePermissionByUserId failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	p := new(models.SearchTableDataParam)
	y := new(models.SysPost)
	list, err := logic.SelectSysPostList(p, y)
	if err != nil {
		zap.L().Error("logic.SelectSysPostList failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	var roleIds []int
	var postids []int
	if useridP != 0 {
		postids, _ = logic.SelectPostListByUserId(useridP)
		roles2, _ := logic.SelectRolePermissionByUserId(int64(useridP))
		for _, sysroles := range roles2 {
			roleIds = append(roleIds, sysroles.RoleId)
		}
	}
	ResponseSuccess(ctx, gin.H{
		"data":    info.SysUser,
		"roles":   roles,
		"posts":   list,
		"postIds": postids,
		"roleIds": roleIds,
	})
}

func SaveUser(ctx *gin.Context) {
	p := new(models.SysUserParam)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("saveUser.shouldbindjson", zap.Error(err))
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
	err = logic.SaveUser(p, id)
	if err != nil {
		zap.L().Error("logic.SaveUser,failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, "保存用户成功")
}

func UploadUser(ctx *gin.Context) {
	p := new(models.SysUserParam)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("saveUser.shouldbindjson", zap.Error(err))
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
	err = logic.UpdateUser(p, id)
	if err != nil {
		zap.L().Error("logic.updateUser failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, "保存用户成功")
}

func DeleteUserById(ctx *gin.Context) {
	userids := ctx.Param("userIds")
	id, err := strconv.Atoi(userids)
	if err != nil {
		zap.L().Error("strconv,itoa", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	err = logic.DeleteUserById(id)
	if err != nil {
		zap.L().Error("logic.DeleteUserById", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "删除用户成功")
}

func ResetPwd(ctx *gin.Context) {
	p := new(models.SysUser)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("resetpwd invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	err := logic.ResetPwd(p)
	if err != nil {
		zap.L().Error("logic.ResetPwd", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "密码更新成功")
}

func ChangeUserStatus(ctx *gin.Context) {
	p := new(models.SysUser)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("resetpwd invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	err := logic.ChangeUserStatus(p)
	if err != nil {
		zap.L().Error("logic.ResetPwd", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "用户状态更新成功")
}

func GetAuthUserRole(ctx *gin.Context) {
	useridstr := ctx.Param("userId")
	atoi, err := strconv.Atoi(useridstr)
	if err != nil {
		zap.L().Error("strconv.Atoi failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}

	info, err := logic.GetUserInfo(int64(atoi))
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, gin.H{
		"user":  info.SysUser,
		"roles": info.Rolepremission,
	})
}

func PutAuthUser(ctx *gin.Context) {
	p := new(models.AddUserRole)
	if err := ctx.ShouldBindJSON(&p); err != nil {
		zap.L().Error("PutAuthUser invalid param failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(ctx, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	err := logic.PutAuthUser(p)
	if err != nil {
		zap.L().Error("logic.PutAuthUser failed", zap.Error(err))
		return
	}
	ResponseSuccess(ctx, "添加成功")
}

// GetUserDeptTree 登录获取菜单
func GetUserDeptTree(ctx *gin.Context) {
	tree, err := logic.GetUserDeptTree()
	if err != nil {
		zap.L().Error("logic.GetUserDeptTree failed", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, tree)
}
