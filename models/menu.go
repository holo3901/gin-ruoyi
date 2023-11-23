package models

import "time"

type SysMenu struct {
	MenuId     int       `db:"menu_id" json:"menuId"` //表示主键
	MenuName   string    `db:"menu_name" json:"menuName"`
	ParentId   int       `db:"parent_id" json:"parentId"`
	OrderNum   int       `db:"order_num" json:"orderNum"`
	MenuType   string    `db:"menu_type" json:"menuType"`
	Visible    string    `db:"visible" json:"visible"`
	Perms      string    `db:"perms" json:"perms"`
	Query      string    `db:"query" json:"query"`
	IsFrame    string    `db:"is_frame" json:"isFrame"`
	Icon       string    `db:"icon" json:"icon"`
	Path       string    `db:"path" json:"path"`
	Status     string    `db:"status" json:"status"`
	IsCache    string    `db:"is_cache" json:"isCache"`
	Component  string    `db:"component" json:"component"`
	CreateBy   string    `db:"create_by" json:"createBy"`
	CreateTime time.Time `db:"create_time" json:"createTime"`
	UpdateBy   string    `db:"update_by" json:"updateBy"`
	UpdateTime time.Time `db:"update_time" json:"updateTime"`
	Remark     string    `db:"remark" json:"remark"`
}

type SysMenuType struct {
	MenuId     int       `db:"menu_id" json:"menuId"` //表示主键
	MenuName   string    `db:"menu_name" json:"menuName"`
	ParentId   int       `db:"parent_id" json:"parentId"`
	OrderNum   int       `db:"order_num" json:"orderNum"`
	MenuType   string    `db:"menu_type" json:"menuType"`
	Visible    string    `db:"visible" json:"visible"`
	Perms      string    `db:"perms" json:"perms"`
	Query      string    `db:"query" json:"query"`
	IsFrame    string    `db:"is_frame" json:"isFrame"`
	Icon       string    `db:"icon" json:"icon"`
	Path       string    `db:"path" json:"path"`
	Status     string    `db:"status" json:"status"`
	IsCache    string    `db:"is_cache" json:"isCache"`
	Component  string    `db:"component" json:"component"`
	CreateBy   string    `db:"create_by" json:"createBy"`
	CreateTime time.Time `db:"create_time" json:"createTime"`
	UpdateBy   string    `db:"update_by" json:"updateBy"`
	UpdateTime time.Time `db:"update_time" json:"updateTime"`
	Remark     string    `db:"remark" json:"remark"`
	RoleId     string    `db:"role_id" json:"roleId"`
}

type MenuVo struct {
	Name       string   `json:"name"`
	Path       string   `json:"path,omitempty"`
	Hidden     bool     `json:"hidden" `
	Redirect   string   `json:"redirect,omitempty"`
	Component  string   `json:"component,omitempty" `
	Query      string   `json:"query,omitempty"`
	AlwaysShow bool     `json:"alwaysShow,omitempty" `
	MetaVo     MetaVo   `json:"meta" `
	Children   []MenuVo `json:"children,omitempty"`
}

type MetaVo struct {
	Title   string `json:"title"`
	Icon    string `json:"icon" `
	NoCache bool   `json:"noCache" `
	Link    string `json:"link,omitempty" `
}

type MenuTreeSelect struct {
	Id       int               `json:"id"`
	Label    string            `json:"label"`
	Children []*MenuTreeSelect `json:"children,omitempty"`
}
