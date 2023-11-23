package models

import "time"

type SysDept struct {
	DeptId     int       `db:"dept_id" json:"deptId"` //表示主键
	ParentId   int       `db:"parent_id" json:"parentId"`
	Ancestors  string    `db:"ancestors" json:"ancestors"`
	DeptName   string    `db:"dept_name" json:"deptName"`
	OrderNum   int       `db:"order_num" json:"orderNum"`
	Leader     string    `db:"leader" json:"leader"`
	Phone      string    `db:"phone" json:"phone"`
	Email      string    `db:"email" json:"email"`
	Status     string    `db:"status" json:"status"`
	DelFlag    string    `db:"del_flag" json:"delFlag"`
	CreateBy   string    `db:"create_by" json:"createBy"`
	CreateTime time.Time `db:"create_time" json:"createTime"`
	UpdateBy   string    `db:"update_by" json:"updateBy"`
	UpdateTime time.Time `db:"update_time" json:"updateTime"`
}

type SysDeptResult struct {
	DeptId     int             `json:"deptId" db:"dept_id"` //表示主键
	ParentId   int             `json:"parentId" db:"parent_id"`
	Ancestors  string          `json:"ancestors" db:"ancestors"`
	DeptName   string          `json:"deptName" db:"dept_name"`
	OrderNum   int             `json:"orderNum" db:"order_num"`
	Leader     string          `json:"leader" db:"leader"`
	Phone      string          `json:"phone" db:"phone"`
	Email      string          `json:"email" db:"email"`
	Status     string          `json:"status" db:"status"`
	DelFlag    string          `db:"del_flag" json:"delFlag"`
	CreateBy   string          `db:"create_by" json:"createBy"`
	CreateTime time.Time       `db:"create_time" json:"createTime"`
	UpdateBy   string          `db:"update_by" json:"updateBy"`
	UpdateTime time.Time       `db:"update_time" json:"updateTime"`
	ParentName string          `json:"parentName"`
	Children   []SysDeptResult `json:"children"`
}

type SysDeptDto struct {
	Id       int           `json:"id"`
	Label    string        `json:"label"`
	Children []*SysDeptDto `json:"children"`
}
