package logic

import (
	"database/sql"
	"errors"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SelectConfigList(param *models.SearchTableDataParam, config *models.SysConfig) ([]*models.SysConfig, error) {
	return mysql.SelectConfigList(param, config)
}

func GetConfigInfo(id int) (*models.SysConfig, error) {
	return mysql.GetConfigInfo(id)
}

func SelectConfig(config *models.SysConfig) (*models.SysConfig, error) {
	return mysql.SelectConfig(config)
}

func SaveConfig(id int64, config *models.SysConfig) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	err = mysql.CheckConfigKeyUnique(config.ConfigKey)
	if err != sql.ErrNoRows {
		return errors.New("key存在")
	}

	config.CreateBy = user.UserName
	config.CreateTime = time.Now()
	return mysql.SaveConfig(config)
}

func UploadConfig(id int64, config *models.SysConfig) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}

	config.UpdateBy = user.UserName
	config.UpdateTime = time.Now()
	return mysql.UploadConfig(config)
}

func DeleteConfig(configid string) error {
	var ids = Split(configid)
	for i := 0; i < len(ids); i++ {
		config, err := mysql.GetConfigInfo(ids[i])
		if err != nil {
			return err
		}
		if config.ConfigType == "Y" {
			return errors.New("内置参数不能删除")
		}
		err = mysql.DeleteConfigById(config.ConfigId)
		if err != nil {
			return err
		}
	}
	return nil
}
