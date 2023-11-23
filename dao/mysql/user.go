package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"ruoyi/models"
	"time"
)

func Login(user *models.User) error {
	oPassword := user.Password
	a := new(models.SysUser)
	sqlStr := `select user_id,user_name,password,status from sys_user where del_flag = '0' and  user_name = ?`
	err := db.Get(a, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}

	if err != nil {
		return err
	}

	if a.Status != "0" {
		return errors.New("账户停用")
	}
	x := PasswordVerify(oPassword, a.Password)

	if x == false {
		return ErrorInvalidPassword
	}
	return nil
}

func FindUserById(id int64) (*models.SysUser, error) {
	a := new(models.SysUser)
	sqlstr := `select user_id,user_name,dept_id,nick_name,user_type,email,phonenumber,sex,avatar,password,create_time,status from sys_user where del_flag = '0' and user_id =? `
	err := db.Get(a, sqlstr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorUserNotExist
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func FindUserByName(name string) (*models.SysUser, error) {
	a := new(models.SysUser)
	sqlstr := `select * from sys_user where del_flag = '0' and user_name =? `
	err := db.Get(a, sqlstr, name)
	if err == sql.ErrNoRows {
		a.UserId = 0
		return a, nil
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}
func GetRolePremissionById(id int64) (roles []*models.SysRoles, err error) {
	a := make([]*models.SysRoles, 0)
	sqlstr := `select distinct r.* from sys_role r left join sys_user_role ur on ur.role_id left join sys_user u on 
               u.user_id = ur.user_id left join sys_dept d on u.dept_id =d.dept_id where r.del_flag = '0' and ur.user_id = ?`
	err = db.Select(&a, sqlstr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorUserNotExist
	}
	return a, nil
}

func GetMenuPermissionById(id int64) (roles []string, err error) {
	a := make([]string, 0)
	sqlstr := `select distinct m.perms from sys_menu m left join sys_role_menu rm on m.menu_id = rm.menu_id 
             left join sys_user_role ur on rm.role_id = ur.role_id 
             left join sys_role r on r.role_id = ur.role_id 
             where m.status = '0' and r.status = '0' and ur.user_id = ?`
	err = db.Select(&a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func SelectUserPostGroup(id int64) ([]string, error) {
	posts := make([]string, 0)
	sqlstr := `select distinct p.post_name from sys_post p 
              left join sys_user_post up on up.post_id = p.post_id
              left join sys_user u on u.user_id =up.user_id
               where u.user_id = ? `
	err := db.Select(&posts, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func UpdateUser(user *models.SysUser) error {
	sqlstr := `update sys_user set nick_name=?,email=?,phonenumber=?,sex=?,password=?,avatar=?,status=? where user_id=?`
	_, err := db.Exec(sqlstr, user.NickName, user.Email, user.Phonenumber, user.Sex, user.Password, user.Avatar, user.Status, user.UserId)
	if err != nil {
		return err
	}
	fmt.Println(err)
	return nil
}

// ----------------------------- 用户管理 ------------------------------

func SelectUserList(param *models.SearchTableDataParam, user *models.SysUser) ([]*models.SysUserParam, error) {
	sqlstr := "SELECT `sys_user`.*, d.dept_name, d.leader FROM `sys_user`" +
		"LEFT JOIN sys_dept d ON d.dept_id = `sys_user`.dept_id WHERE `sys_user`.del_flag = ?"

	var args []interface{}
	args = append(args, "0")

	if user.UserId != 0 {
		sqlstr += " AND `sys_user`.user_id =?"
		args = append(args, user.UserId)
	}
	if user.DeptId != 0 {
		sqlstr += " AND (`sys_user`.dept_id =? OR `sys_user`.dept_id IN(SELECT t.dept_id FROM sys_dept where find_in_set(?,ancestors))) "
		//dept_id在sys_dept表的ancestors字段中包含指定的deptId
		args = append(args, user.DeptId, user.DeptId)
	}
	if user.UserName != "" {
		sqlstr += " AND `sys_user`.user_name like CONCAT('%',?,'%')"
		args = append(args, user.UserName)
	}
	if user.Status != "" {
		sqlstr += " AND `sys_user`.status =?"
		args = append(args, user.Status)
	}
	if user.Phonenumber != "" {
		sqlstr += " AND `sys_user`.phonenumber=?"
		args = append(args, user.Phonenumber)
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += " AND `sys_user`.create_time >=? AND `sys_user`.create_time<=?"
		args = append(args, start, end)
	}

	sqlstr += " order by `sys_user`.user_id "
	if param.PageSize != 0 && param.PageNum != 0 {
		sqlstr += " LIMIT ? OFFSET ?"
		args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	}
	userlist := make([]*models.SysUserParam, 0, 2)
	err := db.Select(&userlist, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return userlist, nil
}

func SelectUserParmList(param *models.SearchTableDataParam, user *models.SysUser) ([]*models.SysUserParam, error) {
	sqlstr := "SELECT `sys_user`.*, d.dept_name, d.leader FROM `sys_user`" +
		"LEFT JOIN sys_dept d ON d.dept_id = `sys_user`.dept_id WHERE `sys_user`.del_flag = ?"

	var args []interface{}
	args = append(args, "0")

	if user.UserId != 0 {
		sqlstr += " AND `sys_user`.user_id =?"
		args = append(args, user.UserId)
	}
	if user.DeptId != 0 {
		sqlstr += " AND (`sys_user`.dept_id =? OR `sys_user`.dept_id IN(SELECT t.dept_id FROM sys_dept where find_in_set(?,ancestors))) "
		//dept_id在sys_dept表的ancestors字段中包含指定的deptId
		args = append(args, user.DeptId, user.DeptId)
	}
	if user.UserName != "" {
		sqlstr += " AND `sys_user`.user_name like CONCAT('%',?,'%')"
		args = append(args, user.UserName)
	}
	if user.Status != "" {
		sqlstr += " AND `sys_user`.status =?"
		args = append(args, user.Status)
	}
	if user.Phonenumber != "" {
		sqlstr += " AND `sys_user`.phonenumber=?"
		args = append(args, user.Phonenumber)
	}
	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += " AND `sys_user`.create_time >=? AND `sys_user`.create_time<=?"
		args = append(args, start, end)
	}
	sqlstr += " LIMIT ? OFFSET ?"
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	userlist := make([]*models.SysUserParam, 0)
	err := db.Select(&userlist, sqlstr, args...)

	if err != nil {
		return nil, err
	}
	return userlist, nil
}

func SaveUser(user *models.SysUserParam) error {
	tx, err := db.Beginx()

	if err != nil {

		return err
	}
	user.LoginDate = time.Now()
	if user.Sex == "男" {
		user.Sex = "0"
	} else if user.Sex == "女" {
		user.Sex = "1"
	} else {
		user.Sex = "2"
	}

	if user.Status == "正常" {
		user.Status = "0"
	} else {
		user.Status = "1"
	}

	sqlstr := `insert into sys_user(nick_name,dept_id,user_name,user_type,
               email,phonenumber,sex,avatar,password,login_date,status,del_flag,create_by,create_time,update_time,remark)
               values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err = tx.Exec(sqlstr, user.NickName, user.DeptId, user.UserName, user.UserType, user.Email,
		user.Phonenumber, user.Sex, user.Avatar, user.Password, user.LoginDate, user.Status,
		user.DelFlag, user.CreateBy, user.CreateTime, user.UpdateTime, user.Remark)
	if err != nil {
		tx.Rollback()
		return err
	}

	sqlstr1 := ` insert into sys_user_post (user_id,post_id) values(?,?)`
	for i := 0; i < len(user.PostIds); i++ {
		_, err = tx.Exec(sqlstr1, user.UserId, user.PostIds[i])
		if err != nil {

			tx.Rollback()
			return err
		}
	}
	sqlstr2 := ` insert into sys_user_role (user_id,role_id) values(?,?)`
	for i := 0; i < len(user.RoleIds); i++ {
		_, err = tx.Exec(sqlstr2, user.UserId, user.RoleIds[i])
		if err != nil {

			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func CheckUserExists(username string) (err error) {
	sqlstr := `select count(user_id) from sys_user where where del_flag=0 AND user_name =?`
	var count int
	if err = db.Get(&count, sqlstr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

func CheckPhoneExists(phonenumber string) (err error) {
	sqlstr := `select count(user_id) from sys_user where where del_flag=0 AND phonenumber =?`
	var count int
	if err = db.Get(&count, sqlstr, phonenumber); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("手机号存在")
	}

	return
}

func CheckEmailExists(email string) (err error) {
	sqlstr := `select count(user_id) from sys_user where del_flag=0 AND email =?`
	var count int
	if err = db.Get(&count, sqlstr, email); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("邮箱存在")
	}

	return
}

func UpdateUserParam(param *models.SysUserParam) (err error) {
	DeleteRoleByUserId(int64(param.UserId))
	AddRoleByUser(param)
	DeletePostByUserId(int64(param.UserId))
	AddPostByUser(param)

	sqlstr := `update sys_user set dept_id=? , user_name=?,nick_name=?,email=?,phonenumber=?,sex=?,status=?,remark=? where user_id=?`
	_, err = db.Exec(sqlstr, param.DeptId, param.UserName, param.NickName, param.Email, param.Phonenumber, param.Sex, param.Status, param.Remark, param.UserId)

	return
}

func DeleteUserById(id int) (err error) {
	sqlstr := `update sys_user set del_flag= '2' where user_id=?`
	_, err = db.Exec(sqlstr, id)

	return
}

func GetUserByDeptId(id int) ([]*models.SysUser, error) {
	sqlstr := `select * from sys_user su
               LEFT JOIN sys_user_role sur ON su.user_id = sur.user_id
               LEFT JOIN sys_role_dept srd ON srd.role_id = sur.role_id 
               where srd.dept_id =? 
               `
	a := make([]*models.SysUser, 0)
	err := db.Select(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}
