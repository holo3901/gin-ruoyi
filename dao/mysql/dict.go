package mysql

import (
	"ruoyi/models"
)

func SelectDictDataList(param *models.SearchTableDataParam, DictData *models.SearchDictData) ([]*models.SysDictData, error) {
	sqlstr := `select * from sys_dict_data where 1=1`

	var args []interface{}

	if DictData.DictLabel != "" {
		sqlstr += ` AND dict_label=?`
		args = append(args, DictData.DictLabel)
	}
	if DictData.DictType != "" {
		sqlstr += ` AND dict_type LIKE ?`
		args = append(args, "%"+DictData.DictType+"%")
	}
	if DictData.Status != "" {
		sqlstr += ` AND status=?`
		args = append(args, DictData.Status)
	}
	sqlstr += ` ORDER BY dict_sort`
	if param.PageSize != 0 && param.PageNum != 0 {
		sqlstr += `  ASC LIMIT ? OFFSET ?`
		args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	}
	post := make([]*models.SysDictData, 0)
	err := db.Select(&post, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func FindDictCodeByType(codetype string) ([]*models.SysDictData, error) {
	sqlstr := `select * from sys_dict_data where status = '0' and dict_type=? order by dict_sort `
	code := make([]*models.SysDictData, 0, 2)
	err := db.Select(&code, sqlstr, codetype)
	if err != nil {
		return nil, err
	}
	return code, nil
}

func FindDictCodeById(dictCode string) (*models.SysDictData, error) {
	sqlstr := `select *from sys_dict_data where dic_code=?`
	var code *models.SysDictData
	err := db.Get(code, sqlstr, dictCode)
	if err != nil {
		return nil, err
	}
	return code, nil
}

func SaveDictData(dict *models.SysDictData) (err error) {
	if dict.IsDefault == "" {
		dict.IsDefault = "N"
	}
	sqlstr := `insert into sys_dict_data(create_by,create_time,update_time,is_default)
              values(?,?,?,?)`
	_, err = db.Exec(sqlstr, dict.CreateBy, dict.CreateTime, dict.UpdateTime, dict.IsDefault)

	return
}

func EditDictData(dict *models.SysDictData) (err error) {
	sqlstr := `update sys_dict_data set dict_sort=?,is_default=?,status=?,update_by=?,update_time=?,remark=? where dict_code=? `
	_, err = db.Exec(sqlstr, dict.DictSort, dict.IsDefault, dict.Status, dict.UpdateBy, dict.UpdateTime, dict.Remark, dict.DictCode)
	return
}

func DeleteDictData(dictcode string) (err error) {
	sql := `delete from sys_dict_data where dict_code=?`
	_, err = db.Exec(sql, dictcode)
	return
}

func SelectDictTypeList(param *models.SearchTableDataParam, search *models.SysDictType) ([]*models.SysDictType, error) {
	sqlstr := `select * from sys_dict_type where 1=1`

	var args []interface{}

	if search.DictName != "" {
		sqlstr += ` AND dict_name like ?`
		args = append(args, "%"+search.DictName+"%")
	}

	if search.Status != "" {
		sqlstr += ` AND status =?`
		args = append(args, search.Status)
	}
	if search.DictType != "" {
		sqlstr += ` AND dict_type like ?`
		args = append(args, "%"+search.DictType+"%")
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND create_time>=? And create_time <=?`
		args = append(args, start, end)
	}

	sqlstr += ` ORDER BY dict_id ASC LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	dictType := make([]*models.SysDictType, 0, 2)
	err := db.Select(&dictType, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return dictType, nil
}

func FindTypeDictById(dictId string) (*models.SysDictType, error) {
	sqlstr := `select * from sys_dict_type where dict_id=?`
	var code *models.SysDictType
	err := db.Get(code, sqlstr, dictId)
	if err != nil {
		return nil, err
	}
	return code, nil
}

func SaveType(dictType *models.SysDictType) (err error) {
	if dictType.Status == "" {
		dictType.Status = "0"
	}
	sqlstr := `insert into sys_dict_type(dict_name,dict_type,status,create_by,create_time,update_time,remark)
               values(?,?,?,?,?,?,?)`
	_, err = db.Exec(sqlstr, dictType.DictName, dictType.DictType, dictType.Status, dictType.CreateBy, dictType.CreateTime, dictType.UpdateTime, dictType.Remark)
	return

}

func UpdateType(dictType *models.SysDictType) (err error) {
	sqlstr := `update sys_dict_type set dict_name=?,dict_type=?,status=?,update_by=?,update_time=?,remark=? where dict_id=? `
	_, err = db.Exec(sqlstr, dictType.DictName, dictType.DictType, dictType.Status, dictType.UpdateBy, dictType.UpdateTime, dictType.Remark, dictType.DictId)
	return
}

func DeleteDataType(id string) (err error) {
	sql := `delete from sys_dict_type where dict_id=?`
	_, err = db.Exec(sql, id)
	return
}
func GetOptionSelect() ([]*models.SysDictType, error) {
	sqlstr := `select * from sys_dict_type`
	dictType := make([]*models.SysDictType, 0, 2)
	err := db.Select(dictType, sqlstr)
	if err != nil {
		return nil, err
	}
	return dictType, nil
}
