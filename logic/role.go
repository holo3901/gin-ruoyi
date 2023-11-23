package logic

import (
	"database/sql"
	"errors"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"strconv"
	"time"
)

func SelectRolePermissionByUserId(id int64) ([]*models.SysRoles, error) {
	return mysql.SelectRolePermissionByUserId(id)
}

func SelectRoleList(param *models.SearchTableDataParam, roles *models.SysRoles) ([]*models.SysRoles, error) {
	return mysql.SelectRoleList(param, roles)
}

func GetRoleInfo(userid int64, roleid int) (*models.SysRoles, error) {
	err := mysql.CheckRoleDataScope(roleid, int(userid))
	if err != nil {
		return nil, err
	}
	id, err := mysql.FindRoleInfoByRoleId(roleid)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func SaveRole(id int64, param *models.SysRolesParam) error {
	user, _ := mysql.FindUserById(id)
	err := mysql.CheckRoleDataScope(param.RoleId, int(id))
	if err != nil {
		return err
	}
	roles, err := mysql.FindRoleByRoleName(param.RoleName)
	if err != nil {
		return err
	}
	if len(roles) != 0 {
		return errors.New("用户存在")
	}
	roles2, err := mysql.FindRoleByRoleKey(param.RoleKey)
	if err != nil {
		return err
	}
	if len(roles2) != 0 {
		return errors.New("用户权限存在")
	}
	//数据范围 1.全部数据权限，2.自定数据权限 3.本部门数据权限，4.本部门及以下数据权限
	param.CreateBy = user.UserName
	param.CreateTime = time.Now()
	param.DataScope = "2"
	param.DelFlag = "0"
	return mysql.SaveRole(param)

}

func UploadRole(id int64, param *models.SysRolesParam) error {
	if param.RoleId == 1 {
		return errors.New("不允许操作超级管理员")
	}
	err := mysql.CheckRoleDataScope(param.RoleId, int(id))
	if err != nil {
		return err
	}
	roles, err := mysql.FindRoleByRoleName(param.RoleName)
	if err != nil {
		return err
	}
	if len(roles) != 0 {
		return errors.New("用户存在")
	}
	roles2, err := mysql.FindRoleByRoleKey(param.RoleKey)
	if err != nil {
		return err
	}
	if len(roles2) != 0 {
		return errors.New("用户权限存在")
	}
	user, _ := mysql.FindUserById(id)
	param.UpdateBy = user.UserName
	param.UpdateTime = time.Now()

	err = mysql.UpdateRole(param)
	if err != nil {
		return err
	}
	//更新role_menu
	err = mysql.DeleteRoleMenuByRoleId(param.RoleId)
	if err != nil {
		return err
	}
	err = mysql.InsertRoleMenu(param)
	if err != nil {
		return err
	}

	//更新role_dept
	err = mysql.DeleteRoleDeptByRole(param.RoleId)
	if err != nil {
		return err
	}
	return mysql.InsertRoleDept(param)
}

func DeleteRole(id int64, roleid string) error {
	roles := Split(roleid)
	for _, v := range roles {
		if v == 1 {
			return errors.New("超级管理不能删除员")
		}
		err := mysql.CheckRoleDataScope(v, int(id))
		if err != nil {
			return err
		}
		_, err = mysql.FindUserRoleByRoleId(v)
		if err != sql.ErrNoRows {
			return err
		}

	}
	err := mysql.DeleteRoleMenu(roleid)
	if err != nil {
		return err
	}
	return mysql.DeleteRoleDept(roleid)
}

func GetRoleOptionSelect() ([]*models.SysRoles, error) {
	return mysql.GetRoleOptionSelect()
}

func GetAllocatedList(param *models.SearchTableDataParam, userParam *models.SysUserParam) ([]*models.SysUser, error) {
	return mysql.GetAllocatedList(param, userParam)
}

func GetUnAllocatedList(param *models.SearchTableDataParam, userParam *models.SysUserParam) ([]*models.SysUser, error) {
	return mysql.GetUnAllocatedList(param, userParam)
}
func CancelRole(param *models.SysUserRolesParam) error {
	return mysql.CancelRole(param)
}

func CancelRoleAll(userid string, roleid int) error {
	return mysql.CancelRoleAll(userid, roleid)

}

func SelectRoleAll(roleid int, userid string, users int64) error {
	err := mysql.CheckRoleDataScope(roleid, int(users))
	if err != nil {
		return err
	}
	var uid = Split(userid)
	for i := 0; i < len(uid); i++ {
		x := &models.SysUserRoles{UserId: uid[i], RoleId: roleid}
		err = mysql.AddUserRole(x)
		if err != nil {
			return err
		}
	}
	return err
}

func GetDeptTreeRole(roleid string) ([]int, error) {
	roleId, _ := strconv.Atoi(roleid)
	role, err := mysql.FindRoleInfoByRoleId(roleId)
	if err != nil {
		return nil, err
	}
	tree, err := mysql.GetDeptTree(role)
	if err != nil {
		return nil, err
	}
	var a []int
	for _, v := range tree {
		a = append(a, v.DeptId)
	}
	return a, nil
}
