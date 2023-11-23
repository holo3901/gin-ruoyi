package models

import "time"

type SysLogininfor struct {
	InfoId        int       `json:"infoId" db:"info_id"` //表示主键
	UserName      string    `json:"userName" db:"user_name"`
	Ipaddr        string    `json:"ipaddr" db:"ipaddr"`
	LoginLocation string    `json:"loginLocation" db:"login_location"`
	Browser       string    `json:"browser" db:"browser"`
	Os            string    `json:"os" db:"os"`
	Status        string    `json:"status" db:"status"`
	Msg           string    `json:"msg" db:"msg"`
	LoginTime     time.Time `json:"loginTime" db:"login_time"`
}

type SysOnline struct {
	OnlineId      int       `json:"OnlineId" db:"online_id"` //表示主键
	UserName      string    `json:"userName" db:"user_name"`
	Ipaddr        string    `json:"ipaddr" db:"ipaddr"`
	Token         string    `json:"token" db:"token"`
	LoginLocation string    `json:"loginLocation" db:"login_location"`
	Browser       string    `json:"browser" db:"browser"`
	Os            string    `json:"os" db:"os"`
	Status        string    `json:"status" db:"status"`
	Msg           string    `json:"msg" db:"msg"`
	LoginTime     time.Time `json:"loginTime" db:"login_time"`
}
