package mysql

import (
	"errors"
	"fmt"
	"ruoyi/models"
	"strings"
)

func SelectMenuTreeByUserId(id int64) ([]*models.SysMenu, error) {
	menu := make([]*models.SysMenu, 0)
	if id == 1 { //1是管理员
		sqlstr := `SELECT * FROM sys_menu WHERE menu_type IN ('M', 'C') AND status = '0'
               ORDER BY parent_id, order_num`
		err := db.Select(&menu, sqlstr)
		if err != nil {
			return nil, err
		}
		return menu, nil

	} else {
		sqlstr := `SELECT DISTINCT m.* FROM sys_menu m 
                 LEFT JOIN sys_role_menu rm ON m.menu_id = rm.menu_id 
                 LEFT JOIN sys_user_role ur ON rm.role_id = ur.role_id 
                 LEFT JOIN sys_role ro ON ur.role_id = ro.role_id 
                 LEFT JOIN sys_user u ON ur.user_id = u.user_id 
                 WHERE u.user_id = ? AND m.menu_type IN ('M', 'C') 
                 AND m.status = 0 AND ro.status = 0 
                 ORDER BY m.parent_id, m.order_num`
		err := db.Select(&menu, sqlstr, id)
		if err != nil {
			return nil, err
		}
		return menu, nil

	}
}

func BuildMenus(list []*models.SysMenu) []*models.MenuVo {
	menuVos := make([]*models.MenuVo, 0)
	for i := 0; i < len(list); i++ {
		var menu = list[i]
		MenuId := menu.MenuId
		parentId := menu.ParentId
		if 0 == parentId {
			path := ""
			if IsInnerLink(menu.Path) {
				path = menu.Path
			}
			var menuVo = &models.MenuVo{
				Hidden: "1" == menu.Visible,
				Query:  menu.Query,
				MetaVo: models.MetaVo{
					Title:   menu.MenuName,
					Icon:    menu.Icon,
					NoCache: "1" == menu.IsCache,
					Link:    path,
				},
				Name:      getRouteName(menu),
				Path:      getRoutePath(menu),
				Component: getComponent(menu),
			}
			if "M" == menu.MenuType {
				if !IsInnerLink(menu.Path) {
					menuVo.AlwaysShow = true
					menuVo.Redirect = "noRedirect"
					menuVo.Children = BuildChildMenus(MenuId, list)
				}
			}
			menuVos = append(menuVos, menuVo)
		}
	}
	return menuVos
}

func BuildChildMenus(ParentId int, lists []*models.SysMenu) []models.MenuVo {
	List := make([]*models.MenuVo, 0)
	for i := 0; i < len(lists); i++ {
		var menu = lists[i]
		var menuId = menu.MenuId
		var pId = menu.ParentId
		if pId == ParentId {
			var path = ""
			if IsInnerLink(menu.Path) {
				path = menu.Path
			}
			var menuVo = &models.MenuVo{
				Hidden: "1" == menu.Visible,
				Query:  menu.Query,
				MetaVo: models.MetaVo{
					Title:   menu.MenuName,
					Icon:    menu.Icon,
					NoCache: "1" == menu.IsCache,
					Link:    path,
				},
				Name:      getRouteName(menu),
				Path:      getRoutePath(menu),
				Component: getComponent(menu),
			}
			if "M" == menu.MenuType {
				if !IsInnerLink(menu.Path) {
					menuVo.AlwaysShow = true
					menuVo.Redirect = "noRedirect"
					menuVo.Children = BuildChildMenus(menuId, lists)
				}
			}
			List = append(List, menuVo)
		}
	}
	var result []models.MenuVo
	for _, ptr := range List {
		if ptr != nil {
			result = append(result, *ptr)
		}
	}

	return result
}

func getRouteName(menu *models.SysMenu) string {
	var name = FirstUpper(menu.Path)
	if isMenuFrame(menu) {
		return ""
	}
	return name
}

func getComponent(menu *models.SysMenu) string {
	var component = "Layout"
	if "" != menu.Component && !isMenuFrame(menu) {
		component = menu.Component
	} else if "" == menu.Component && IsInnerLink(menu.Path) {
		component = "InnerLink"
	} else if "" == menu.Component && isParentView(menu) {
		component = "parentView"
	}
	return component
}
func getRoutePath(menu *models.SysMenu) string {
	var routerPath = menu.Path
	if IsInnerLink(routerPath) {
		return routerPath
	}
	if 0 == menu.ParentId && "M" == menu.MenuType && "1" == menu.IsFrame {
		routerPath = "/" + menu.Path
	} else if isMenuFrame(menu) {
		routerPath = "/"
	}
	return routerPath
}

func isParentView(menu *models.SysMenu) bool {
	return menu.ParentId != 0 && "M" == menu.MenuType
}

// 是否为外链
func IsInnerLink(path string) bool {
	return strings.Contains(path, "http://") || strings.Contains(path, "https://")
}

func isMenuFrame(menu *models.SysMenu) bool {
	return menu.ParentId == 0 && "C" == menu.MenuType && menu.IsFrame == "1"
}

/*首字母大写*/
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func SelectSysMenuList(menu *models.SysMenu, param *models.SearchTableDataParam) ([]*models.SysMenu, error) {
	var sql = "select menu_id, menu_name, parent_id, order_num, path, component, `query`, is_frame, is_cache, menu_type, visible, " +
		"status, ifnull(perms,'') as perms, icon, create_time from sys_menu where 1=1"
	var name = menu.MenuName
	var args []interface{}
	if name != "" {
		sql += " AND menu_name like ?"
		args = append(args, "%"+name+"%")
	}
	var visible = menu.Visible
	if visible != "" {
		sql += " AND visible = ?"
		args = append(args, visible)
	}
	var status = menu.Status
	if status != "" {
		sql += " AND status = ?"
		args = append(args, status)
	}
	rows := make([]*models.SysMenu, 0)
	err := db.Select(&rows, sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func SelectSysMenuListByUserId(id int64, menu *models.SysMenu, param *models.SearchTableDataParam) ([]*models.SysMenu, error) {
	var sql = "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible," +
		" m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time " +
		"from sys_menu m left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"left join sys_user_role ur on rm.role_id = ur.role_id " +
		"left join sys_role ro on ur.role_id = ro.role_id where ur.user_id = ? "
	var args []interface{}

	args = append(args, id)

	if menu.MenuName != "" {
		sql += " AND m.menu_name like ?"
		args = append(args, "%"+menu.MenuName+"%")
	}

	if menu.Visible != "" {
		sql += " AND m.visible = ?"
		args = append(args, menu.Visible)
	}

	if menu.Status != "" {
		sql += " AND m.status = ?"
		args = append(args, menu.Status)
	}

	sql += " order by m.parent_id,m.order_num"
	list := make([]*models.SysMenu, 0)
	err := db.Select(&list, sql, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetMenuInfo(menu int) (*models.SysMenu, error) {
	sqlstr := "select menu_id,menu_name,parent_id,order_num,path,component," +
		"`query`,is_frame,is_cache,menu_type,visible,status,ifnull(perms,'') as perms,icon" +
		",create_time from sys_menu where menu_id =?"
	list := new(models.SysMenu)
	err := db.Get(list, sqlstr, menu)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func SelectMenuTree(id int64, menu *models.SysMenu) ([]*models.SysMenu, error) {
	menus := make([]*models.SysMenu, 0)
	sql := "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time from sys_menu m"
	var args []interface{}
	if id == 1 {
		if menu.MenuName != "" {
			sql += " AND m.menu_name like ?"
			args = append(args, "%"+menu.MenuName+"%")
		}
		if menu.Visible != "" {
			sql += " AND m.visible = ?"
			args = append(args, menu.Visible)
		}
		if menu.Status != "" {
			sql += " AND m.status =?"
			args = append(args, menu.Status)
		}
		err := db.Select(&menus, sql, args...)
		if err != nil {
			return nil, err
		}
		return menus, nil
	}

	sql += " left join sys_role_menu rm on m.menu_id = rm.menu_id" +
		"left join sys_user_role ur on rm.role_id = ur.role_id" +
		"left koin sys_role ro on ur.role_id =ro.role_id" +
		"where ur.user_id =?"
	args = append(args, id)
	if menu.MenuName != "" {
		sql += " AND m.menu_name like ?"
		args = append(args, "%"+menu.MenuName+"%")
	}
	if menu.Visible != "" {
		sql += " AND m.visible = ?"
		args = append(args, menu.Visible)
	}
	if menu.Status != "" {
		sql += " AND m.status =?"
		args = append(args, menu.Status)
	}
	err := db.Select(&menus, sql, args...)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

func BuildMenusTree(menu []*models.SysMenu) []*models.MenuTreeSelect {
	menus := make([]*models.MenuTreeSelect, 0)
	for i := 0; i < len(menu); i++ {
		menuId := menu[i].MenuId
		parentId := menu[i].ParentId
		if parentId == 0 {
			var menVo = &models.MenuTreeSelect{
				Id:    menuId,
				Label: menu[i].MenuName,
			}
			menVo.Children = BuildChildMenuTreeSelect(menuId, menu)
			menus = append(menus, menVo)
		}
	}
	return menus
}

func BuildChildMenuTreeSelect(parentId int, lists []*models.SysMenu) []*models.MenuTreeSelect {
	var list []*models.MenuTreeSelect
	for i := 0; i < len(lists); i++ {
		var menuId = lists[i].MenuId
		var pId = lists[i].ParentId
		if pId == parentId {
			var menVo = &models.MenuTreeSelect{
				Id:    menuId,
				Label: lists[i].MenuName,
			}
			menVo.Children = BuildChildMenuTreeSelect(menuId, lists)
			list = append(list, menVo)
		}
	}
	return list
}

func SelectMenuTreeByRoleId(id int, roles *models.SysRoles) ([]*models.SysMenuType, error) {
	sqlstr := "select * from sys_menu m left join sys_role_menu rm on m.menu_id =rm.menu_id " +
		"where rm.role_id = ?"
	var args []interface{}
	args = append(args, id)
	if roles.MenuCheckStrictly {
		sqlstr += " AND m.menu_id not in ( select m.parent_id from sys_menu m inner join sys_role_menu rm on  m.menu_id = rm.menu_id where rm.role_id = ? ) "
		args = append(args, id)
	}
	sqlstr += " order by m.parent_id,m.order_num"
	menu := make([]*models.SysMenuType, 0)
	fmt.Println("sqlstr", sqlstr)
	fmt.Println("args", args)
	err := db.Select(&menu, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	fmt.Println(menu)
	return menu, nil
}

func CheckMenuNameUnique(parent int, menuname string) error {
	sql := `select count(*) from sys_menu where menu_name = ? AND parent_id =?`
	var count int
	err := db.Get(&count, sql, menuname, parent)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("菜单名存在")
	}
	return nil
}

func AddMenu(menu *models.SysMenu) error {
	sqlstr := "insert into sys_menu(menu_name,parent_id,order_num,path,component,`query`,is_frame," +
		"is_cache,menu_type,perms,icon,create_by,create_time,remark)" +
		"values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := db.Exec(sqlstr, menu.MenuName, menu.ParentId, menu.OrderNum, menu.Path, menu.Component, menu.Query, menu.IsFrame, menu.IsCache,
		menu.MenuType, menu.Perms, menu.Icon, menu.CreateBy, menu.CreateTime, menu.Remark)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMenu(menu *models.SysMenu) error {
	sqlstr := "update sys_menu set menu_name=?,parent_id=?,order_num=?,path=?,component=?,is_frame=?,is_cache=?,menu_type=?,perms=?,icon=?,update_by=?,update_time=?,remark=? where menu_id=? "
	_, err := db.Exec(sqlstr, menu.MenuName, menu.ParentId, menu.OrderNum, menu.Path, menu.Component, menu.IsFrame, menu.IsCache, menu.MenuType, menu.Perms, menu.Icon, menu.UpdateBy, menu.UpdateTime, menu.Remark, menu.MenuId)
	if err != nil {
		return err
	}
	return nil
}

func HasChildrenByMenuId(id int) ([]*models.SysMenu, error) {
	sqlstr := `select * from sys_menu where parent_id =?`
	a := make([]*models.SysMenu, 0)
	err := db.Select(&a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func CheckMenuExistRole(id int) error {
	sql := `select count(*) from sys_role_menu where menu_id =?`
	var count int
	err := db.Get(&count, sql, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("菜单已分配，无法删除")
	}
	return nil
}

func DeleteMenu(id int) error {
	sqlstr := `delete from sys_menu where menu_id =?`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}
