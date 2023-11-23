package mysql

import (
	"github.com/gin-gonic/gin"
	useragent "github.com/wenlng/go-user-agent"
	"ruoyi/models"
	ipinfo "ruoyi/pkg/ip"
	"time"
)

func LoginInfoAdd(ctx *gin.Context, param *models.LoginParam, user *models.User, message string, loginsuccess bool) {
	var status = "0"
	if loginsuccess {
		status = "0"
	} else {
		status = "1"
	}
	userAgent := ctx.Request.Header.Get("user-Agent")
	Os := useragent.GetOsName(userAgent)
	browser := useragent.GetBrowserName(userAgent)
	ip := ipinfo.GetRemoteClientIp(ctx.Request)
	var info = &models.SysLogininfor{
		UserName:      param.Username,
		Ipaddr:        "" + ip,
		LoginLocation: "" + ipinfo.GetRealAddressByIP(ip),
		Browser:       "" + browser,
		Os:            "" + Os,
		Status:        status,
		Msg:           message,
		LoginTime:     time.Now(),
	}
	AddOnline(info)

	sqlstr := `insert into sys_logininfor(user_name,ipaddr,login_location,browser,os,status,msg,login_time)
               values (?,?,?,?,?,?,?,?)`
	_, _ = db.Exec(sqlstr, info.UserName, info.Ipaddr, info.LoginLocation, info.Browser, info.Os, info.Status, info.Msg, info.LoginTime)
}

func DeleteLoginInfo(id int64) {
	byId, err := FindUserById(id)
	if err != nil {
		return
	}
	sqlstr := `delete from sys_logininfor where user_name =?`
	db.Exec(sqlstr, byId.UserName)
}

func SelectLogininforList(param *models.SearchTableDataParam, logininfor *models.SysLogininfor) ([]*models.SysLogininfor, error) {
	sqlstr := `select * from sys_logininfor where 1=1`
	var args []interface{}
	if logininfor.Ipaddr != "" {
		sqlstr += ` AND ipaddr like ?`
		args = append(args, logininfor.Ipaddr+"%")
	}
	if logininfor.Status != "" {
		sqlstr += ` AND status=?`
		args = append(args, logininfor.Status)
	}
	if logininfor.UserName != "" {
		sqlstr += ` AND user_name like ?`
		args = append(args, logininfor.UserName+"%")
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND login_time >=? AND login_time <=?`
		args = append(args, start, end)
	}

	if param.OrderByColumn != "" {
		if param.IsAsc == "ascending" {
			if param.OrderByColumn == "loginTime" {
				sqlstr += ` order by login_time DESC`
			}
			if param.OrderByColumn == "userName" {
				sqlstr += ` order by user_name DESC`
			}
		}
		if param.IsAsc == "descending" {
			if param.OrderByColumn == "loginTime" {
				sqlstr += ` order by login_time ASC`
			}
			if param.OrderByColumn == "userName" {
				sqlstr += ` order by user_name ASC`
			}
		}
	}
	sqlstr += "order by info_id DESC LIMIT ? OFFSET ? "
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	list := make([]*models.SysLogininfor, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func DeleteInfoId(id string) error {
	sqlstr := `delete from sys_logininfor where id in (?)`
	_, err := db.Exec(sqlstr, id)
	return err
}

func ClearLoginLog() error {
	sql := `TRUNCATE TABLE sys_logininfor`
	_, err := db.Exec(sql)
	return err
}
