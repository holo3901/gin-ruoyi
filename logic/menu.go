package logic

import (
	"errors"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SelectSysMenuByUserId(id int64, menu *models.SysMenu, param *models.SearchTableDataParam) ([]*models.SysMenu, error) {
	if id == 1 {
		list, err := mysql.SelectSysMenuList(menu, param)
		if err != nil {
			return nil, err
		}
		return list, nil
	}

	userId, err := mysql.SelectSysMenuListByUserId(id, menu, param)
	if err != nil {
		return nil, err
	}
	return userId, nil
}

func GetMenuInfo(id int) (*models.SysMenu, error) {
	return mysql.GetMenuInfo(id)
}

func GetTreeSelect(id int64, menu *models.SysMenu) ([]*models.MenuTreeSelect, error) {
	tree, err := mysql.SelectMenuTree(id, menu)
	if err != nil {
		return nil, err
	}
	return mysql.BuildMenusTree(tree), nil
}

func GetMenuTreeSelectByRole(id int) ([]int, error) {
	roleId, err := mysql.FindRoleInfoByRoleId(id)
	if err != nil {
		return nil, err
	}
	byRoleId, err := mysql.SelectMenuTreeByRoleId(id, roleId)
	if err != nil {
		return nil, err
	}
	var a []int
	for _, v := range byRoleId {
		a = append(a, v.MenuId)
	}
	return a, nil
}

func SaveMenu(id int64, menu *models.SysMenu) error {
	byId, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	if err = mysql.CheckMenuNameUnique(menu.ParentId, menu.MenuName); err != nil {
		return err
	}
	isPath := mysql.IsInnerLink(menu.Path)
	if menu.IsFrame == "true" && !isPath {
		return errors.New("新增菜单失败，地址必须以https开头")
	}
	menu.CreateBy = byId.UserName
	menu.CreateTime = time.Now()

	return mysql.AddMenu(menu)
}

func UpdateMenu(id int64, menu *models.SysMenu) error {
	byId, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	if menu.ParentId == menu.MenuId {
		return errors.New("修改菜单失败，上级菜单不能为自己")
	}
	if err = mysql.CheckMenuNameUnique(menu.ParentId, menu.MenuName); err != nil {
		return err
	}
	isPath := mysql.IsInnerLink(menu.Path)
	if menu.IsFrame == "true" && !isPath {
		return errors.New("修改菜单失败，地址必须以https开头")
	}

	menu.UpdateBy = byId.UpdateBy
	menu.UpdateTime = time.Now()
	return mysql.UpdateMenu(menu)
}

func DeleteMenu(id int) error {
	menuId, err := mysql.HasChildrenByMenuId(id)
	if err != nil {
		return err
	}
	if len(menuId) > 0 {
		return errors.New("存在子菜单，无法删除")
	}
	if err = mysql.CheckMenuExistRole(id); err != nil {
		return err
	}
	return mysql.DeleteMenu(id)
}
