package models

import "time"

// SysRoles model：数据库字段
type SysRoles struct {
	RoleId            int       `json:"roleId" db:"role_id"` //表示主键
	RoleName          string    `json:"roleName" db:"role_name"`
	RoleKey           string    `json:"roleKey" db:"role_key"`
	RoleSort          int       `json:"roleSort" db:"role_sort"`
	DataScope         string    `json:"dataScope" db:"data_scope"`
	Status            string    `json:"status" db:"status"`
	MenuCheckStrictly bool      `json:"menuCheckStrictly" db:"menu_check_strictly"`
	DeptCheckStrictly bool      `json:"deptCheckStrictly" db:"dept_check_strictly"`
	DelFlag           string    `json:"delFlag" db:"del_flag"`
	CreateBy          string    `json:"createBy" db:"create_by"`
	CreateTime        time.Time `json:"create_time" db:"create_time"`
	UpdateBy          string    `json:"updateBy" db:"update_by"`
	UpdateTime        time.Time `json:"updateTime" db:"update_time"`
	Remark            string    `json:"remark" db:"remark"`
}

type SysRolesParam struct {
	RoleId            int       `json:"roleId" db:"role_id"` //表示主键
	RoleName          string    `json:"roleName" db:"role_name"`
	RoleKey           string    `json:"roleKey" db:"role_key"`
	RoleSort          int       `json:"roleSort" db:"role_sort"`
	DataScope         string    `json:"dataScope" db:"data_scope"`
	Status            string    `json:"status" db:"status"`
	MenuCheckStrictly bool      `json:"menuCheckStrictly" db:"menu_check_strictly"`
	DeptCheckStrictly bool      `json:"deptCheckStrictly" db:"dept_check_strictly"`
	DelFlag           string    `json:"delFlag" db:"del_flag"`
	CreateBy          string    `json:"createBy" db:"create_by"`
	CreateTime        time.Time `json:"createTime" db:"create_time"`
	UpdateBy          string    `json:"updateBy" db:"update_by"`
	UpdateTime        time.Time `json:"updateTime" db:"update_time"`
	Remark            string    `json:"remark" db:"remark"`
	MenuIds           []int     `json:"menuIds"`
	DeptIds           []int     `json:"deptIds"`
}
