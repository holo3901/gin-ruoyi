package mysql

import "ruoyi/models"

func GetDeptList(param *models.SearchTableDataParam, dept *models.SysDept) ([]*models.SysDeptResult, error) {
	sqlstr := `select * from sys_dept where 1=1`

	var args []interface{}
	if dept.DeptId != 0 {
		sqlstr += ` AND dept_id =?`
		args = append(args, dept.DeptId)
	}

	if dept.ParentId != 0 {
		sqlstr += ` AND parent_id =?`
		args = append(args, dept.ParentId)
	}
	if dept.DeptName != "" {
		sqlstr += ` AND dept_name like ?`
		args = append(args, "%"+dept.DeptName+"%")
	}
	if dept.Status != "" {
		sqlstr += ` AND status=?`
		args = append(args, dept.Status)
	}
	sqlstr += ` ORDER BY parent_id, order_num ASC `
	if param.PageNum != 0 && param.PageSize != 0 {
		sqlstr += ` limit ? OFFSET ?`
		args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	}
	deptlist := make([]*models.SysDeptResult, 0, 2)
	err := db.Select(&deptlist, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return deptlist, nil
}

func GetDeptInfo(id int) (*models.SysDept, error) {
	dept := new(models.SysDept)
	sqlstr := `select * from  sys_dept where dept_id= ?`
	err := db.Get(dept, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func SaveDept(dept *models.SysDept) error {
	sqlstr := `insert into sys_dept (parent_id,ancestors,dept_name,order_num,leader,phone,email,status,create_by,create_time)
               values (?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlstr, dept.ParentId, dept.Ancestors, dept.DeptName, dept.OrderNum, dept.Leader, dept.Phone, dept.Email, dept.Status, dept.CreateBy, dept.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

func UploadDept(dept *models.SysDept) error {
	sqlstr := `update sys_dept set parent_id=?,ancestors=?,dept_name=?,order_num=?,leader=?,phone=?,email=?,status=?,update_by=?,update_time=? where dept_id=? `
	_, err := db.Exec(sqlstr, dept.ParentId, dept.Ancestors, dept.DeptName, dept.OrderNum, dept.Leader, dept.Phone, dept.Email, dept.Status, dept.UpdateBy, dept.UpdateTime, dept.DeptId)
	if err != nil {
		return err
	}
	return nil
}

func GetDeptByParentId(id int) ([]*models.SysDept, error) {
	sqlstr := `select * from sys_dept where parent_id =?`
	a := make([]*models.SysDept, 0)
	err := db.Select(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func DeleteDept(id int) error {
	sqlstr := `delete from sys_dept where dept_id =?`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}
