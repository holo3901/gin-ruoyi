package models

import "time"

// SysDictType model：数据库字段
type SysDictType struct {
	DictId     int       `json:"dictId" db:"dict_id"` //表示主键
	DictName   string    `json:"dictName" db:"dict_name"`
	DictType   string    `json:"dictType" db:"dict_type"`
	Status     string    `json:"status" db:"status"`
	CreateBy   string    `json:"createBy" db:"create_by"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	UpdateBy   string    `json:"updateBy" db:"update_by"`
	UpdateTime time.Time `json:"updateTime" db:"update_time"`
	Remark     string    `json:"remark" db:"remark"`
}
