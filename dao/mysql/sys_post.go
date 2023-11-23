package mysql

import "ruoyi/models"

func SelectSysPostList(param *models.SearchTableDataParam, post *models.SysPost) ([]*models.SysPost, error) {
	sqlstr := `select * from sys_post where 1=1`
	var args []interface{}
	if post.PostCode != "" {
		sqlstr += ` AND post_code like ?`
		args = append(args, "%"+post.PostCode+"%")
	}

	if post.Status != "" {
		sqlstr += ` AND status=?`
		args = append(args, post.Status)
	}

	if post.PostName != "" {
		sqlstr += `AND post_name like ?`
		args = append(args, post.PostName)
	}
	sqlstr += ` order by post_sort`

	if param.PageSize != 0 {
		sqlstr += ` LIMIT ? OFFSET ?`
		args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	}
	posts := make([]*models.SysPost, 0, 2)
	err := db.Select(&posts, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func SelectPostListByUserId(id int) ([]*models.SysPost, error) {
	a := make([]*models.SysPost, 0)
	sqlstr := `select * from sys_post p left join sys_user_post on up.post_id = p.post_d
              left join sys_user u on u.user_id = up.user_id where u.user_id=?`

	err := db.Select(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func DeletePostByUserId(id int64) (err error) {
	sqlstr := `delete from sys_user_post where user_id =?`
	_, err = db.Exec(sqlstr, id)
	return
}

func AddPostByUser(param *models.SysUserParam) (err error) {
	sqlstr := `insert into sys_user_post (user_id,role_id) values (?,?)`
	for _, v := range param.PostIds {
		_, err = db.Exec(sqlstr, param.UserId, v)
		if err != nil {
			return err
		}
	}
	return
}

func FindPostInfoById(id int) (*models.SysPost, error) {
	sqlstr := `select * from sys_post where post_id =?`
	post := new(models.SysPost)
	err := db.Get(post, sqlstr)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func SavePost(post *models.SysPost) error {
	sqlstr := `insert into sys_post (post_code,post_name,post_sort,create_by,create_time,remark) values(?,?,?,?,?,?)`
	_, err := db.Exec(sqlstr, post.PostCode, post.PostName, post.PostSort, post.CreateBy, post.CreateTime, post.Remark)
	if err != nil {
		return err
	}
	return nil
}

func UploadPost(post *models.SysPost) error {
	sqlstr := `update sys_post set post_code=?,post_name=?,post_sort=?,update_by=?,update_time=?,remark=? where post_id =?`
	_, err := db.Exec(sqlstr, post.PostCode, post.PostName, post.PostSort, post.UpdateBy, post.UpdateTime, post.Remark, post.PostId)
	if err != nil {
		return err
	}
	return nil
}

func DeletePost(id int) error {
	sqlstr := `delete from sys_post where post_id =?`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}
