package logic

import (
	"ruoyi/dao/mysql"
	"ruoyi/dao/redis"
	"ruoyi/models"
)

func ListOnline(param *models.SearchTableDataParam, online *models.SysOnline) ([]*models.SysOnline, error) {
	return mysql.ListOnline(param, online)
}

func DeleteOnline(id int) error {
	byId, err := mysql.FindOnlineById(id)
	if err != nil {
		return err
	}
	name, err := mysql.FindUserByName(byId.UserName)
	if err != nil {
		return err
	}
	err = redis.CreateLogin(int64(name.UserId), "")
	if err != nil {
		return err
	}
	return mysql.DeleteOnline(id)
}
