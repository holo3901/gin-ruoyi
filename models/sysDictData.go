package models

import "time"

// SysDictData model：数据库字段
type SysDictData struct {
	DictCode   int       `json:"dictCode" db:"dict_code"` //表示主键
	DictSort   int       `json:"dictSort" db:"dict_sort"`
	DictLabel  string    `json:"dictLabel" db:"dict_label"`
	DictValue  string    `json:"dictValue" db:"dict_value"`
	DictType   string    `json:"dictType" db:"dict_type"`
	CssClass   string    `db:"css_class" json:"cssClass"`
	ListClass  string    `db:"list_class" json:"listClass"`
	IsDefault  string    `db:"is_default" json:"isDefault"`
	Status     string    `db:"status" json:"status"`
	CreateBy   string    `db:"create_by" json:"createBy"`
	CreateTime time.Time `db:"create_time" json:"createTime"`
	UpdateBy   string    `db:"update_by" json:"updateBy"`
	UpdateTime time.Time `db:"update_time" json:"updateTime"`
	Remark     string    `db:"remark" json:"remark"`
}
