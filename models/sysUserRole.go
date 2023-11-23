package models

type SysUserRoles struct {
	RoleId int `json:"roleId"`
	UserId int `json:"userId"`
}

type SysUserRolesParam struct {
	RoleId string `json:"roleId"`
	UserId int    `json:"userId"`
}
