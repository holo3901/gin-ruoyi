package models

import "time"

// SysUser model：数据库字段
type SysUser struct {
	UserId      int       `json:"user_id" db:"user_id"` //表示主键
	DeptId      int       `json:"dept_id" db:"dept_id"`
	UserName    string    `json:"user_name" db:"user_name"`
	NickName    string    `json:"nick_name" db:"nick_name"`
	UserType    string    `json:"user_type" db:"user_type"`
	Email       string    `json:"email" db:"email"`
	Phonenumber string    `json:"phonenumber" db:"phonenumber"`
	Sex         string    `json:"sex" db:"sex"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Password    string    `json:"password" db:"password"`
	Status      string    `json:"status" db:"status"`
	DelFlag     string    `json:"del_flag" db:"del_flag"`
	LoginIp     string    `json:"login_ip" db:"login_ip"`
	LoginDate   time.Time `json:"login_date" db:"login_date"`
	CreateBy    string    `json:"create_by" db:"create_by"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateBy    string    `json:"update_by" db:"update_by"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
	Remark      string    `json:"remark" db:"remark"`
}

type SysUserExcel struct {
	UserId      int       `json:"user_id" db:"user_id"` //表示主键
	DeptId      int       `json:"dept_id" db:"dept_id"`
	UserName    string    `json:"user_name" db:"user_name"`
	NickName    string    `json:"nick_name" db:"nick_name"`
	UserType    string    `json:"user_type" db:"user_type"`
	Email       string    `json:"email" db:"email"`
	Phonenumber string    `json:"phonenumber" db:"phonenumber"`
	Sex         string    `json:"sex" db:"sex"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Password    string    `json:"password" db:"password"`
	Status      string    `json:"status" db:"status"`
	DelFlag     string    `json:"del_flag" db:"del_flag"`
	LoginIp     string    `json:"login_ip" db:"login_ip"`
	LoginDate   time.Time `json:"login_date" db:"login_date"`
	CreateBy    string    `json:"create_by" db:"create_by"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateBy    string    `json:"update_by" db:"update_by"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
	Remark      string    `json:"remark" db:"remark"`
	DeptName    string    `json:"deptName" db:"dept_name"`
	Leader      string    `json:"leader" db:"leader"`
}

type SysUserParam struct {
	UserId      int       `json:"user_id" db:"user_id"` //表示主键
	DeptId      int       `json:"dept_id" db:"dept_id"`
	UserName    string    `json:"user_name" db:"user_name"`
	NickName    string    `json:"nick_name" db:"nick_name"`
	UserType    string    `json:"user_type" db:"user_type"`
	Email       string    `json:"email" db:"email"`
	Phonenumber string    `json:"phonenumber" db:"phonenumber"`
	Sex         string    `json:"sex" db:"sex"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Password    string    `json:"password" db:"password"`
	Status      string    `json:"status" db:"status"`
	DelFlag     string    `json:"del_flag" db:"del_flag"`
	LoginIp     string    `json:"login_ip" db:"login_ip"`
	LoginDate   time.Time `json:"login_date" db:"login_date"`
	CreateBy    string    `json:"create_by" db:"create_by"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateBy    string    `json:"update_by" db:"update_by"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
	Remark      string    `json:"remark" db:"remark"`
	DeptName    string    `json:"deptName" db:"dept_name"`
	Leader      string    `json:"leader" db:"leader"`
	PostIds     []int     `json:"postIds"`
	RoleIds     []int     `json:"roleIds"`
	RoleId      int       `json:"roleId"`
}
