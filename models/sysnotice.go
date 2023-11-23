package models

import "time"

type SysNotice struct {
	NoticeId      int       `json:"noticeId" db:"notice_id""` //表示主键
	NoticeTitle   string    `json:"noticeTitle" db:"notice_title"`
	NoticeType    string    `json:"noticeType" db:"notice_type"`
	NoticeContent string    `json:"noticeContent" db:"notice_content"`
	Status        string    `json:"status" db:"status"`
	CreateBy      string    `json:"createBy" db:"create_by"`
	CreateTime    time.Time `json:"createTime" db:"create_time"`
	UpdateBy      string    `json:"updateBy" db:"update_by"`
	UpdateTime    time.Time `json:"updateTime" db:"update_time"`
	Remark        string    `json:"remark" db:"remark"`
}
