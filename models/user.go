package models

import (
	"encoding/json"
)

// login里对用户的注册
type User struct {
	UserID   int64  `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"user_name"`
	Password string `json:"password" db:"password"`
	Token    string
}

type LoginParam struct {
	Code     string `json:"code"`
	Password string `json:"password"`
	Username string `json:"username"`
	Uuid     string `json:"uuid"`
}

type Repassword struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,nefield=OldPassword"`
}
type UserInfo struct {
	PostGroup      []string
	Rolepremission []*SysRoles
	MenuPremission []string
	*SysUser
	*SysDept
}

// Userparam 请求参数
type Userparam struct {
	NickName    string `json:"nickName"`
	UserName    string `json:"userName"`
	Phonenumber string `json:"phonenumber"`
	Email       string `json:"email" binding:"email" msg:"邮箱地址格式不正确"`
	Sex         string `json:"sex"`
}

type SearchTableDataParam struct {
	PageNum       int             `json:"pageNum"`
	PageSize      int             `json:"pageSize"`
	Other         json.RawMessage `json:"other"`
	OrderByColumn string          `json:"orderByColumn"`
	IsAsc         string          `json:"isAsc"`
	Params        Params          `json:"params"`
}

type Params struct {
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
}

type SearchDictData struct {
	DictType  string `json:"dict_type"`
	Status    string `json:"status"`
	DictLabel string `json:"dict_label"`
}
type AddUserRole struct {
	Userid  int    `json:"user_id"`
	RoleIds string `json:"role_ids"`
}
