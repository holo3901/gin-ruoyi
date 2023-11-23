package models

import "time"

type SysConfig struct {
	ConfigId    int       `json:"configId" db:"config_id"` //表示主键
	ConfigName  string    `json:"configName" db:"config_name"`
	ConfigKey   string    `json:"configKey" db:"config_key"`
	ConfigValue string    `json:"configValue" db:"config_value"`
	ConfigType  string    `json:"configType" db:"config_type"`
	CreateBy    string    `json:"createBy" db:"create_by"`
	CreateTime  time.Time `json:"createTime" db:"create_time"`
	UpdateBy    string    `json:"updateBy" db:"update_by"`
	UpdateTime  time.Time `json:"updateTime" db:"update_time"`
	Remark      string    `json:"remark" db:"remark"`
}
