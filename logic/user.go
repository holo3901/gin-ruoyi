package logic

import (
	"errors"
	"fmt"
	"mime/multipart"
	"ruoyi/dao/mysql"
	"ruoyi/dao/redis"
	"ruoyi/models"
	"ruoyi/pkg/JWT"
	"ruoyi/pkg/qiniu"
	"strconv"
	"time"
)

func Login(p *models.LoginParam) (*models.User, error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	users, _ := mysql.FindUserByName(user.Username)

	user.UserID = int64(users.UserId)
	token, err := JWT.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}

	user.Token = token
	err = redis.CreateLogin(user.UserID, user.Token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserInfo(id int64) (*models.UserInfo, error) {
	sys, err := mysql.FindUserById(id)
	if err != nil {
		return nil, err
	}
	Rolepremission, err := mysql.GetRolePremissionById(id)
	if err != nil {
		return nil, err
	}

	MenuPremission, err := mysql.GetMenuPermissionById(id)
	if err != nil {
		return nil, err
	}
	Deptinfo, err := mysql.GetDeptInfo(sys.DeptId)
	if err != nil {
		return nil, err
	}
	Postgroup, err := mysql.SelectUserPostGroup(id)
	if err != nil {
		return nil, err
	}
	return &models.UserInfo{
		PostGroup:      Postgroup,
		Rolepremission: Rolepremission,
		MenuPremission: MenuPremission,
		SysUser:        sys,
		SysDept:        Deptinfo,
	}, nil
}

func GetRouters(id int64) ([]*models.MenuVo, error) {
	_, err := mysql.FindUserById(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(id)
	menu, err := mysql.SelectMenuTreeByUserId(id)

	return mysql.BuildMenus(menu), nil
}

func UpdatePwd(id int64, oldpwd string, newpwd string) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	if !mysql.PasswordVerify(oldpwd, user.Password) {
		return errors.New("密码错误")
	}
	hash, err := mysql.PasswordHash(newpwd)
	if err != nil {
		return errors.New("加密失败")
	}
	user.Password = hash
	err = mysql.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func EditProfileUserInfo(id int64, info *models.Userparam) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	if info.NickName != "" {
		user.NickName = info.NickName
	}
	if info.Email != "" {
		user.Email = info.Email
	}
	if info.Phonenumber != "" {
		user.Phonenumber = info.Phonenumber
	}
	if info.Sex != "" {
		user.Sex = info.Sex
	}
	err = mysql.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAvatar(file multipart.File, filesize int64, name string, id int64) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return errors.New("用户不存在")
	}
	avator, err := qiniu.UploadToQiNiu(file, filesize, name)
	if err != nil {
		return err
	}
	user.Avatar = avator
	err = mysql.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

// ----------------------------- 用户管理 ------------------------------

func SelectUserList(param *models.SearchTableDataParam, user *models.SysUser) ([]*models.SysUserParam, error) {
	return mysql.SelectUserList(param, user)
}
func SelectUserParmList(param *models.SearchTableDataParam, user *models.SysUser) ([]*models.SysUserParam, error) {
	return mysql.SelectUserParmList(param, user)
}

func ImportUserData(users []*models.SysUserParam, updateSupport string, id int64) (string, error) {
	var errList []string
	var errorsNum int
	user, _ := mysql.FindUserById(id)
	for i := 0; i < len(users); i++ {

		users[i].CreateTime = time.Now()
		var u2, _ = mysql.FindUserByName(users[i].UserName)
		password, err := mysql.SelectConfigByKey("sys.user.initPassword")
		if err != nil {
			return "", err
		}
		if password != "" {
			pwd, _ := mysql.PasswordHash(password)
			users[i].Password = pwd
		} else {
			pwd, _ := mysql.PasswordHash("123456")
			users[i].Password = pwd
		}

		if u2.UserId == 0 {
			users[i].CreateBy = user.UserName
			users[i].CreateTime = time.Now()
		}
		var d = strconv.Itoa(i + 1)
		users[i].UpdateBy = user.UserName
		users[i].UpdateTime = time.Now()
		users[i].DelFlag = "0"
		if u2.UserId == 0 || updateSupport == "true" { //用户添加
			err := mysql.SaveUser(users[i])

			if err != nil {
				errorsNum += 1
				errList = append(errList, d+",用户名"+users[i].UserName+"，添加失败<br/>")
			} else {
				errList = append(errList, d+",用户名"+users[i].UserName+"，添加成功")
			}

		} else { //用户更新
			errorsNum += 1
			errList = append(errList, d+"、用户名："+users[i].UserName+"，已存在<br/>")

		}
	}
	var result string
	for _, v := range errList {
		result += v
	}
	if errorsNum >= 1 {
		return result, errors.New("数据添加失败")
	}
	return result, nil
}

func CheckUserDataScope(userid int64, useridp string) error {
	ID, _ := strconv.Atoi(useridp)
	param := new(models.SearchTableDataParam)
	x := new(models.SysUser)
	x.UserId = ID
	list, _ := mysql.SelectUserList(param, x)
	if len(list) == 0 {
		return errors.New("没有数据")
	}
	return nil
}
func SaveUser(user *models.SysUserParam, id int64) error {
	if err := mysql.CheckUserExists(user.UserName); err != nil {
		return err
	}
	if err := mysql.CheckPhoneExists(user.Phonenumber); err != nil {
		return err
	}
	if err := mysql.CheckEmailExists(user.Email); err != nil {
		return err
	}
	pwd, _ := mysql.PasswordHash(user.Password)
	user1, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}

	user.CreateBy = user1.UserName
	user.CreateTime = time.Now()

	user.Password = pwd
	err = mysql.SaveUser(user)
	if err != nil {
		return err
	}
	return err
}

func UpdateUser(param *models.SysUserParam, id int64) (err error) {
	user1, err := mysql.FindUserById(id)
	user2, err := mysql.FindUserById(int64(param.UserId))
	if err != nil {
		return err
	}
	if err = mysql.CheckUserExists(param.UserName); err != nil {
		if user1.UserName != user2.UserName {
			return err
		}
	}
	if err = mysql.CheckPhoneExists(param.Phonenumber); err != nil {
		if user1.Phonenumber != user2.Phonenumber {
			return err
		}
	}
	if err = mysql.CheckEmailExists(param.Email); err != nil {
		if user1.Email != user2.Email {
			return err
		}
	}

	param.UpdateBy = user1.UserName
	param.UpdateTime = time.Now()
	err = mysql.UpdateUserParam(param)
	if err != nil {
		return err
	}

	return
}

func DeleteUserById(id int) (err error) {
	if id == 1 {
		return errors.New("不能对管理员进行操作")
	}
	err = mysql.DeletePostByUserId(int64(id))
	if err != nil {
		return err
	}
	err = mysql.DeleteRoleByUserId(int64(id))
	if err != nil {
		return err
	}
	return mysql.DeleteUserById(id)
}

func ResetPwd(param *models.SysUser) (err error) {
	user, err := mysql.FindUserById(int64(param.UserId))
	if err != nil {
		return err
	}
	if param.Password != "" {
		var pwd, _ = mysql.PasswordHash(param.Password)
		user.Password = pwd
	}
	err = mysql.UpdateUser(user)
	return
}

func ChangeUserStatus(param *models.SysUser) (err error) {
	user, err := mysql.FindUserById(int64(param.UserId))
	if err != nil {
		return err
	}
	user.Status = param.Status
	err = mysql.UpdateUser(user)
	return
}

func PutAuthUser(role *models.AddUserRole) (err error) {
	roles := new(models.SysUserParam)
	data := Split(role.RoleIds)
	roles.UserId = role.Userid
	roles.RoleIds = data
	err = mysql.AddRoleByUser(roles)
	return
}

func GetUserDeptTree() ([]*models.SysDeptDto, error) {
	param := new(models.SearchTableDataParam)
	sysdept := new(models.SysDept)
	list, err := mysql.GetDeptList(param, sysdept)
	if err != nil {
		return nil, err
	}
	data := make([]*models.SysDeptDto, 0)
	for i := 0; i < len(list); i++ {
		var bean = list[i]
		if bean.ParentId == 0 {
			dept := &models.SysDeptDto{
				Id:    bean.DeptId,
				Label: bean.DeptName,
			}
			data = append(data, getDeptChildren(list, dept))
		}
	}
	return data, nil
}

func getDeptChildren(list []*models.SysDeptResult, dept *models.SysDeptDto) *models.SysDeptDto {
	var data []*models.SysDeptDto
	for i := 0; i < len(list); i++ {
		var bean = list[i]
		if bean.ParentId == dept.Id {
			depts := &models.SysDeptDto{
				Id:    bean.DeptId,
				Label: bean.DeptName,
			}
			data = append(data, getDeptChildren(list, depts))
		}
	}
	dept.Children = data
	return dept
}
