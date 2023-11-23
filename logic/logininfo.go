package logic

import (
	"ruoyi/dao/mysql"
	"ruoyi/models"
)

func LoginInformList(param *models.SearchTableDataParam, logininfor *models.SysLogininfor) ([]*models.SysLogininfor, error) {
	return mysql.SelectLogininforList(param, logininfor)
}

func DeleteInfoId(id string) error {
	return mysql.DeleteInfoId(id)
}

func ClearLoginLog() error {
	return mysql.ClearLoginLog()
}
