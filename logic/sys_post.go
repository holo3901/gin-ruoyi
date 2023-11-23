package logic

import (
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SelectSysPostList(param *models.SearchTableDataParam, syspost *models.SysPost) ([]*models.SysPost, error) {
	return mysql.SelectSysPostList(param, syspost)
}

func SelectPostListByUserId(id int) ([]int, error) {
	post, err := mysql.SelectPostListByUserId(id)
	if err != nil {
		return nil, err
	}
	var a []int
	for i := 0; i < len(post); i++ {
		a = append(a, post[i].PostId)
	}
	return a, nil
}

func FindPostInfoById(id int) (*models.SysPost, error) {
	return mysql.FindPostInfoById(id)
}

func SavePost(id int64, post *models.SysPost) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	post.CreateBy = user.UserName
	post.CreateTime = time.Now()
	return mysql.SavePost(post)
}

func UploadPost(id int64, post *models.SysPost) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	post.UpdateBy = user.UserName
	post.UpdateTime = time.Now()
	return mysql.UploadPost(post)
}

func DeletePost(id int) error {
	return mysql.DeletePost(id)
}
