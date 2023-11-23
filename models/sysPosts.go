package models

import "time"

// SysPost model：数据库字段
type SysPost struct {
	PostId     int       `json:"postId" db:"post_id"` //表示主键
	PostCode   string    `json:"postCode" db:"post_code"`
	PostName   string    `json:"postName" db:"post_name"`
	PostSort   int       `json:"postSort" db:"post_sort"`
	Status     string    `json:"status" db:"status"`
	CreateBy   string    `json:"createBy" db:"create_by"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
	UpdateBy   string    `json:"updateBy" db:"update_by"`
	UpdateTime time.Time `json:"updateTime" db:"update_time"`
	Remark     string    `json:"remark" db:"remark"`
}
