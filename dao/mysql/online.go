package mysql

import (
	"ruoyi/models"
)

func AddOnline(info *models.SysLogininfor) {
	sqlstr := `insert into sys_online(user_name,ipaddr,login_location,browser,os,status,msg,login_time)
               values (?,?,?,?,?,?,?,?)`
	_, _ = db.Exec(sqlstr, info.UserName, info.Ipaddr, info.LoginLocation, info.Browser, info.Os, info.Status, info.Msg, info.LoginTime)
}

func DeleteOnline(id int) error {
	sqlstr := `delete from sys_online where online_id =?`
	_, err := db.Exec(sqlstr, id)
	return err
}

func DeleteOnlineByName(name string) error {
	sqlstr := `delete from sys_online where user_name =?`
	_, err := db.Exec(sqlstr, name)
	return err
}

func DeleteOnlineAll() {
	sqlstr := `delete from sys_online`
	db.Exec(sqlstr)
}

func ListOnline(param *models.SearchTableDataParam, online *models.SysOnline) ([]*models.SysOnline, error) {
	sqlstr := `select * from sys_online where 1=1`
	var args []interface{}
	if online.Ipaddr != "" {
		sqlstr += ` AND ipaddr like ?`
		args = append(args, "%"+online.Ipaddr+"%")
	}
	if online.UserName != "" {
		sqlstr += ` AND user_name like ?`
		args = append(args, "%"+online.UserName+"%")
	}
	sqlstr += ` LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	list := make([]*models.SysOnline, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindOnlineById(id int) (*models.SysOnline, error) {
	sqlstr := `select * from sys_online where online_id =?`
	a := new(models.SysOnline)
	err := db.Get(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}
