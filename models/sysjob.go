package models

import "time"

type SysJobParam struct {
	JobId          int
	Concurrent     int
	CronExpression string
	InvokeTarget   string
	JobGroup       string
	JobName        string
	MisfirePolicy  string
	Status         string
}
type SysJob struct {
	JobId          int       `json:"jobId" db:"job_id"` //表示主键
	JobName        string    `json:"jobName" db:"job_name"`
	JobGroup       string    `json:"jobGroup" db:"job_group"`
	InvokeTarget   string    `json:"invokeTarget" db:"invoke_target"`
	CronExpression string    `json:"cronExpression" db:"cron_expression"`
	MisfirePolicy  int       `json:"misfirePolicy" db:"misfire_policy"`
	Concurrent     int       `json:"concurrent" db:"concurrent"`
	Status         string    `json:"status" db:"status"`
	CreateBy       string    `json:"createBy" db:"create_by"`
	CreateTime     time.Time `json:"createTime" db:"create_time"`
	UpdateBy       string    `json:"updateBy" db:"update_by"`
	UpdateTime     time.Time `json:"updateTime" db:"update_time"`
	Remark         string    `json:"remark" db:"remark"`
}

type SysJobLog struct {
	JobLogId      int       `json:"jobLogId" db:"job_log_id"` //表示主键
	JobName       string    `json:"jobName" db:"job_name"`
	JobGroup      string    `json:"jobGroup" db:"job_group"`
	InvokeTarget  string    `json:"invokeTarget" db:"invoke_target"`
	JobMessage    string    `json:"jobMessage" db:"job_message"`
	ExceptionInfo string    `json:"exceptionInfo" db:"exception_info"`
	StartTime     string    `json:"startTime"`
	StopTime      string    `json:"stopTime"`
	Status        string    `json:"status" db:"status"`
	CreateTime    time.Time `json:"createTime" db:"create_time"`
}
