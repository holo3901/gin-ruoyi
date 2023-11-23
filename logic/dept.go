package logic

import (
	"errors"
	"github.com/wxnacy/wgo/arrays"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"strconv"
	"time"
)

func GetDeptList(param *models.SearchTableDataParam, dept *models.SysDept) ([]*models.SysDeptResult, error) {
	return mysql.GetDeptList(param, dept)
}

func ExcludeDept(id string, param *models.SearchTableDataParam, dept *models.SysDept) ([]*models.SysDeptResult, error) {
	list, err := mysql.GetDeptList(param, dept)
	if err != nil {
		return nil, err
	}
	excludeList := make([]*models.SysDeptResult, 0)
	for i := 0; i < len(list); i++ {
		answer := SplitStr(list[i].Ancestors)
		index := arrays.ContainsString(answer, id) //包含字符串 返回字符串 val 在数组中的索引位置
		if id != strconv.Itoa(list[i].DeptId) || index == -1 {
			excludeList = append(excludeList, list[i])
		}
	}
	return excludeList, nil
}

func GetDeptInfo(id int) (*models.SysDept, error) {
	return mysql.GetDeptInfo(id)
}

func SaveDept(dept *models.SysDept, id int64) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	deptinfo, err := mysql.GetDeptInfo(dept.ParentId)
	if err != nil {
		return err
	}
	dept.Ancestors = deptinfo.DeptName + strconv.Itoa(dept.ParentId)
	dept.CreateBy = user.UserName
	dept.CreateTime = time.Now()
	return mysql.SaveDept(dept)
}

func UploadDept(dept *models.SysDept, id int64) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	deptinfo, err := mysql.GetDeptInfo(dept.ParentId)
	if err != nil {
		return err
	}
	dept.Ancestors = deptinfo.DeptName + strconv.Itoa(dept.ParentId)
	dept.UpdateBy = user.UserName
	dept.UpdateTime = time.Now()
	return mysql.UploadDept(dept)
}

func DeleteDept(id int) error {
	parentId, err := mysql.GetDeptByParentId(id)
	if err != nil {
		return err
	}
	if len(parentId) > 0 {
		return errors.New("存在下级部门，无法删除")
	}

	deptId, err := mysql.GetUserByDeptId(id)
	if err != nil {
		return err
	}
	if len(deptId) > 0 {
		return errors.New("存在部门用户，无法删除")
	}
	return mysql.DeleteDept(id)
}
