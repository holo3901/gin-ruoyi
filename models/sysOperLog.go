package models

import "time"

// SysOperLog model：数据库字段
type SysOperLog struct {
	OperId        int       `json:"operId" db:"oper_id"` //表示主键
	Title         string    `json:"title" db:"title"`
	BusinessType  string    `json:"businessType" db:"business_type"`
	Method        string    `json:"method" db:"method"`
	RequestMethod string    `json:"requestMethod" db:"request_method"`
	OperatorType  string    `json:"operatorType" db:"operator_type"`
	OperName      string    `json:"operName" db:"oper_name"`
	DeptName      string    `json:"deptName" db:"dept_name"`
	OperUrl       string    `json:"operUrl" db:"oper_url"`
	OperIp        string    `json:"operIp" db:"oper_ip"`
	OperLocation  string    `json:"operLocation" db:"oper_location"`
	OperParam     string    `json:"operParam" db:"oper_param"`
	JsonResult    string    `json:"jsonResult" db:"json_result"`
	Status        string    `json:"status" db:"status"`
	ErrorMsg      string    `json:"errorMsg" db:"error_msg"`
	OperTime      time.Time `json:"operTime" db:"oper_time"`
}
