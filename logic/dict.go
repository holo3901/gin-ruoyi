package logic

import (
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SelectDictDataList(param *models.SearchTableDataParam, DictData *models.SearchDictData) ([]*models.SysDictData, error) {
	list, err := mysql.SelectDictDataList(param, DictData)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindDictCodeByID(dictcode string) (*models.SysDictData, error) {
	return mysql.FindDictCodeById(dictcode)
}
func FindDictCodeByType(dictType string) ([]*models.SysDictData, error) {
	return mysql.FindDictCodeByType(dictType)
}

func SaveDictData(id int64, dictDataParam *models.SysDictData) (err error) {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return
	}
	if dictDataParam.DictCode == 0 {
		dictDataParam.CreateBy = user.UserName
		dictDataParam.CreateTime = time.Now()
		dictDataParam.UpdateTime = time.Now()

		return mysql.SaveDictData(dictDataParam)
	}
	dictDataParam.UpdateBy = user.UserName
	dictDataParam.UpdateTime = time.Now()
	return mysql.EditDictData(dictDataParam)
}

func DeleteDictData(dictcode string) (err error) {
	return mysql.DeleteDictData(dictcode)
}

func SelectSysDictTypeList(param *models.SearchTableDataParam, search *models.SysDictType) ([]*models.SysDictType, error) {

	list, err := mysql.SelectDictTypeList(param, search)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindTypeDictById(dictId string) (*models.SysDictType, error) {
	return mysql.FindTypeDictById(dictId)
}

func SaveType(dictType *models.SysDictType, id int64) error {
	byId, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	if dictType.DictId == 0 {
		dictType.CreateBy = byId.UserName
		dictType.CreateTime = time.Now()
		dictType.UpdateTime = time.Now()
		return mysql.SaveType(dictType)
	}
	dictType.UpdateBy = byId.UserName
	dictType.UpdateTime = time.Now()
	return mysql.UpdateType(dictType)
}

func DeleteDataType(id string) error {
	return mysql.DeleteDataType(id)
}

func GetOptionSelect() ([]*models.SysDictType, error) {
	return mysql.GetOptionSelect()
}
