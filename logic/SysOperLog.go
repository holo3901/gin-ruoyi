package logic

import (
	"ruoyi/dao/mysql"
	"ruoyi/models"
)

func GetOperLogList(param *models.SearchTableDataParam, log *models.SysOperLog) ([]*models.SysOperLog, error) {
	return mysql.GetOperLogList(param, log)
}

func DeleteOperlog(id int) error {
	return mysql.DeleteOperlog(id)
}

func ClearOperlog() error {
	return mysql.ClearOperlog()
}
