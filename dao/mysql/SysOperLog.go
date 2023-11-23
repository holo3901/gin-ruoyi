package mysql

import (
	"ruoyi/models"
)

func GetOperLogList(param *models.SearchTableDataParam, log *models.SysOperLog) ([]*models.SysOperLog, error) {
	sqlstr := `select * from sys_oper_log where 1=1`
	var args []interface{}

	if log.OperIp != "" {
		sqlstr += ` AND oper_ip like ?`
		args = append(args, "%"+log.OperIp+"%")
	}
	if log.Title != "" {
		sqlstr += ` AND title like ?`
		args = append(args, "%"+log.Title+"%")
	}
	if log.OperName != "" {
		sqlstr += ` AND oper_name like ? `
		args = append(args, "%"+log.OperName+"%")
	}
	if log.RequestMethod != "" {
		sqlstr += ` AND request_method like ?`
		args = append(args, "%"+log.RequestMethod+"%")
	}
	if log.Status != "" {
		sqlstr += ` AND status =?`
		args = append(args, log.Status)
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND start_time>=? AND start_time <=?`
		args = append(args, start, end)
	}
	list := make([]*models.SysOperLog, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func DeleteOperlog(id int) error {
	sqlstr := `delete from sys_oper_log where oper_id =?`
	_, err := db.Exec(sqlstr, id)
	return err
}

func ClearOperlog() error {
	sqlstr := `truncate table sys_oper_log`
	_, err := db.Exec(sqlstr)
	return err
}
