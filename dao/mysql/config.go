package mysql

import (
	"fmt"
	"ruoyi/models"
)

func SelectConfigByKey(configkey string) (string, error) {
	sqlstr := `select * from sys_config where config_key like ?`
	config := new(models.SysConfig)
	err := db.Get(config, sqlstr, "%"+configkey+"%")
	if err != nil {
		return "", err
	}

	return config.ConfigValue, nil
}

func SelectConfigList(param *models.SearchTableDataParam, config *models.SysConfig) ([]*models.SysConfig, error) {
	sqlstr := `select * from sys_config where 1=1`
	var args []interface{}
	if config.ConfigId != 0 {
		sqlstr += ` AND config_id =?`
		args = append(args, config.ConfigId)
	}

	if config.ConfigKey != "" {
		sqlstr += ` AND config_key = ?`
		args = append(args, config.ConfigKey)
	}

	if config.ConfigName != "" {
		sqlstr += ` AND config_name like ?`
		args = append(args, "%"+config.ConfigName+"%")
	}

	if config.ConfigType != "" {
		sqlstr += ` AND config_type = ?`
		args = append(args, config.ConfigType)
	}

	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND create_time >= ? AND create_time <= ?`
		args = append(args, start, end)
	}

	sqlstr += ` LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	fmt.Println(sqlstr)
	fmt.Println(args)
	list := make([]*models.SysConfig, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func GetConfigInfo(id int) (*models.SysConfig, error) {
	sqlstr := `select * from sys_config where config_id = ?`
	a := new(models.SysConfig)
	err := db.Get(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SelectConfig(config *models.SysConfig) (*models.SysConfig, error) {
	sqlstr := `select * from sys_config where 1=1`
	var args []interface{}
	if config.ConfigId != 0 {
		sqlstr += ` AND config_id = ?`
		args = append(args, config.ConfigId)
	}
	if config.ConfigKey != "" {
		sqlstr += ` AND config_key = ?`
		args = append(args, config.ConfigKey)
	}
	if config.ConfigName != "" {
		sqlstr += ` AND config_name like ?`
		args = append(args, config.ConfigName)
	}
	a := new(models.SysConfig)

	err := db.Get(a, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SaveConfig(config *models.SysConfig) error {
	sqlstr := ` insert into sys_config (config_name,config_key,config_value,create_by,create_time,remark) 
                values (?,?,?,?,?,?)`
	_, err := db.Exec(sqlstr, config.ConfigName, config.ConfigKey, config.ConfigValue, config.CreateBy, config.CreateTime, config.Remark)
	if err != nil {
		return err
	}
	return nil
}

func CheckConfigKeyUnique(key string) error {
	sqlstr := `select * from sys_config where config_key =?`
	a := new(models.SysConfig)
	err := db.Get(a, sqlstr, key)
	if err != nil {
		return err
	}
	return nil
}

func UploadConfig(config *models.SysConfig) error {
	sqlstr := ` update sys_config set config_name=? config_key=?,config_value=?,update_by=?,update_time=?,remark=? where config_id =?`
	_, err := db.Exec(sqlstr, config.ConfigName, config.ConfigKey, config.ConfigValue, config.UpdateBy, config.UpdateTime, config.Remark, config.ConfigId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteConfigById(id int) error {
	sqlstr := `delete from sys_config where config_id =?`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}
