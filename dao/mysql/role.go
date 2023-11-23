package mysql

import (
	"errors"
	"ruoyi/models"
)

const baseSql = `select distinct r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly, 
	r.status, r.del_flag, r.create_time, r.remark 
	from sys_role r 
	left join sys_user_role ur on ur.role_id = r.role_id
	left join sys_user u on u.user_id = ur.user_id 
	left join sys_dept d on u.dept_id = d.dept_id `

func SelectRolePermissionByUserId(id int64) ([]*models.SysRoles, error) {
	sqlstr := `select distinct r.* from sys_role r 
               left join sys_user_role ur on ur.role_id =r.role_id
               left join sys_user u on u.user_id =ur.user_id
               left join sys_dept d on u.dept_id = d.dept_id
               where r.del_flag = '0' and ur.user_id=?`
	roles := make([]*models.SysRoles, 0)
	err := db.Select(&roles, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func DeleteRoleByUserId(id int64) (err error) {
	sqlstr := `delete from sys_user_role where user_id =? `
	_, err = db.Exec(sqlstr, id)
	return
}

func AddRoleByUser(param *models.SysUserParam) (err error) {
	sqlstr := `insert into sys_user_role (user_id,role_id) values(?,?) `
	for _, v := range param.RoleIds {
		_, err = db.Exec(sqlstr, param.UserId, v)
		if err != nil {
			return err
		}
	}
	return
}

func FindRoleInfoByRoleId(id int) (*models.SysRoles, error) {
	sqlstr := `select * from sys_role where role_id =? `
	a := new(models.SysRoles)
	err := db.Get(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SelectRoleList(param *models.SearchTableDataParam, roles *models.SysRoles) ([]*models.SysRoles, error) {
	sqlstr := `select distinct sys_role.* FROM sys_role 
		      left join sys_user_role ur on ur.role_id =sys_role.role_id
		      left join sys_user u on u.user_id =ur.user_id
		      left join sys_dept d on u.dept_id =d.dept_id
              where sys_role.del_flag = '0'`
	var args []interface{}
	if roles.RoleId != 0 {
		sqlstr += ` AND sys_role.role_id=?`
		args = append(args, roles.RoleId)
	}
	if roles.RoleName != "" {
		sqlstr += ` AND sys_role.role_name like ?`
		args = append(args, "%"+roles.RoleName+"%")
	}
	if roles.RoleKey != "" {
		sqlstr += ` AND sys_role.role_key like ?`
		args = append(args, "%"+roles.RoleKey+"%")
	}

	if roles.Status != "" {
		sqlstr += ` AND sys_role.status=?`
		args = append(args, roles.Status)
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND sys_role.create_time>=? AND sys_role.create_time<=?`
		args = append(args, start, end)
	}
	sqlstr += `ORDER BY sys_role.role_sort LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)

	list := make([]*models.SysRoles, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// 检查角色是否有数据权限
func CheckRoleDataScope(roleid int, id int) error {
	if id == 1 {
		return nil
	}
	var sql = baseSql + `where r.del_flag= '0' AND r.role_id =? AND u.user_id=? order by r.role_sort`
	a := make([]*models.SysRoles, 0)
	err := db.Select(&a, sql, roleid, id)
	if err != nil {
		return err
	}
	if len(a) < 1 {
		return errors.New("没有权限访问角色数据")
	}
	return nil
}

// 校验角色是否存在
func FindRoleByRoleName(name string) ([]*models.SysRoles, error) {
	sqlstr := `select * from sys_role where role_name =?`
	a := make([]*models.SysRoles, 0)
	err := db.Select(a, sqlstr, name)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// 校验角色权限是否存在
func FindRoleByRoleKey(name string) ([]*models.SysRoles, error) {
	sqlstr := `select * from sys_role where role_key =?`
	a := make([]*models.SysRoles, 0)
	err := db.Select(a, sqlstr, name)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SaveRole(roles *models.SysRolesParam) error {
	sqlstr := `insert into sys_role (role_name,role_key,role_sort,data_scope,menu_check_strictly,dept_check_strictly
               ,status,del_flag,create_by,create_time,remark) values(?,?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlstr, roles.RoleName, roles.RoleKey, roles.RoleSort, roles.DataScope, roles.MenuCheckStrictly, roles.DeptCheckStrictly,
		roles.Status, roles.DelFlag, roles.CreateBy, roles.CreateTime, roles.Remark)
	if err != nil {
		return err
	}
	return nil
}
func UpdateRole(roles *models.SysRolesParam) error {
	sqlstr := `update sys_role set role_name=?,role_sort=?,date_scope=?,menu_check_strictly=?,dept_check_strictly=?
             ,status=?,del_flag=?,update_by=?,update_time=?,remark=? where role_id =?`
	_, err := db.Exec(sqlstr, roles.RoleName, roles.RoleKey, roles.RoleSort, roles.DataScope, roles.MenuCheckStrictly, roles.DeptCheckStrictly,
		roles.Status, roles.DelFlag, roles.UpdateBy, roles.UpdateTime, roles.Remark, roles.RoleId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoleMenuByRoleId(roleid int) error {
	sqlstr := `delete from sys_role_menu where role_id = ?`
	_, err := db.Exec(sqlstr, roleid)
	if err != nil {
		return err
	}
	return nil
}
func DeleteRoleMenu(id string) error {
	sqlstr := `delete from sys_role_menu where role_id in (?)`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}

func InsertRoleMenu(role *models.SysRolesParam) error {
	tx, _ := db.Beginx()
	sqlstr := `insert into sys_role_menu (role_id,menu_id) values(?,?) `

	for _, v := range role.MenuIds {
		_, err := tx.Exec(sqlstr, role.RoleId, v)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err := tx.Commit()
	return err
}

func DeleteRoleDeptByRole(roleid int) error {
	sqlstr := `delete from sys_role_dept where role_id = ?`
	_, err := db.Exec(sqlstr, roleid)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRoleDept(id string) error {
	sqlstr := `delete from sys_role_dept where role_id in (?)`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}
func InsertRoleDept(role *models.SysRolesParam) error {
	tx, _ := db.Beginx()
	sqlstr := `insert into sys_role_dept (role_id,dept_id) values(?,?) `

	for _, v := range role.DeptIds {
		_, err := tx.Exec(sqlstr, role.RoleId, v)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err := tx.Commit()
	return err
}

func FindUserRoleByRoleId(id int) ([]*models.SysUser, error) {
	sql := `select u.* from sys_user u 
    left join sys_user_role ur on ur.user_id =u.user_id
    left join sys_role r on r.role_id =ur.role_id 
    where r.del_flag= '0' AND r.role_id=?`
	a := make([]*models.SysUser, 0)
	err := db.Select(a, sql, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func GetRoleOptionSelect() ([]*models.SysRoles, error) {
	sqlstr := baseSql + `where r.del_flag='0' order by r.role_sort`
	a := make([]*models.SysRoles, 0)
	err := db.Select(a, sqlstr)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func GetAllocatedList(param *models.SearchTableDataParam, userParam *models.SysUserParam) ([]*models.SysUser, error) {
	sql := `select distinct u.*  from sys_user u 
            left join sys_dept d on u.dept_id = d.dept_id
            left join sys_user_role ur on u.user_id=ur.user_id
            left join sys_role r on r.role_id = ur.role_id
            where u.del_flag='0' and r.role_id = ?`
	var args []interface{}
	args = append(args, userParam.RoleId)
	if userParam.UserName != "" {
		sql += ` AND u.user_name like ?`
		args = append(args, "%"+userParam.UserName+"%")
	}
	if userParam.Phonenumber != "" {
		sql += ` AND u.phonenumber like ?`
		args = append(args, "%"+userParam.Phonenumber+"%")
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sql += ` AND u.create_time>=? AND u.create_time<=?`
		args = append(args, start, end)
	}
	sql += ` LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)

	a := make([]*models.SysUser, 0)
	err := db.Select(&a, sql, args...)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func GetUnAllocatedList(param *models.SearchTableDataParam, userParam *models.SysUserParam) ([]*models.SysUser, error) {
	sql := `select distinct u.*  from sys_user u 
            left join sys_dept d on u.dept_id = d.dept_id
            left join sys_user_role ur on u.user_id=ur.user_id
            left join sys_role r on r.role_id = ur.role_id
            where u.del_flag='0' and (r.role_id = ? or r.role_id IS NULL)
            and u.user_id not in (select u.user_id from sys_user u 
            inner join sys_user_role ur on u.user_id = ur.user_id and ur.role_id = ? `
	var args []interface{}
	args = append(args, userParam.RoleId, userParam.RoleId)
	if userParam.UserName != "" {
		sql += ` AND u.user_name like ?`
		args = append(args, "%"+userParam.UserName+"%")
	}
	if userParam.Phonenumber != "" {
		sql += ` AND u.phonenumber like ?`
		args = append(args, "%"+userParam.Phonenumber+"%")
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sql += ` AND u.create_time>=? AND u.create_time<=?`
		args = append(args, start, end)
	}
	sql += ` LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)

	a := make([]*models.SysUser, 0)
	err := db.Select(&a, sql, args...)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func CancelRole(param *models.SysUserRolesParam) error {
	sqlStr := `delete from sys_user_role where user_id = ? AND role_id =?`
	_, err := db.Exec(sqlStr, param.UserId, param.RoleId)
	if err != nil {
		return err
	}
	return nil
}

func CancelRoleAll(userids string, roleid int) error {
	sqlStr := `delete from sys_user_role where user_id in (?) AND role_id =?`
	_, err := db.Exec(sqlStr, userids, roleid)
	if err != nil {
		return err
	}
	return nil
}

func AddUserRole(roles *models.SysUserRoles) error {
	sqlstr := `insert into sys_user_role (user_id,role_id) values (?,?)`
	_, err := db.Exec(sqlstr, roles.UserId, roles.RoleId)
	if err != nil {
		return err
	}
	return nil
}

func GetDeptTree(dept *models.SysRoles) ([]*models.SysDept, error) {
	sqlstr := `select d.* from sys_dept d left join sys_role_dept rd on d.dept_id=rd.dept_id
               where rd.role_id = ?`
	var args []interface{}
	args = append(args, dept.RoleId)
	if dept.DeptCheckStrictly {
		sqlstr += ` and d.dept_id not in (select d.parent_id from sys_dept d inner join sys_role_dept rd on d.dept_id=rd.dept_id and rd.role_id = ? )`
		args = append(args, dept.RoleId)
	}
	sqlstr += ` order by d.parent_id,d.order_num`

	a := make([]*models.SysDept, 0)
	err := db.Select(&a, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return a, nil
}
