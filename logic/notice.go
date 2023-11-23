package logic

import (
	"errors"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SelectSysNoticeList(param *models.SearchTableDataParam, notice *models.SysNotice) ([]*models.SysNotice, error) {
	return mysql.SelectSysNoticeList(param, notice)
}

func FindNoticeInfoById(id int) (*models.SysNotice, error) {
	return mysql.FindNoticeInfoById(id)
}

func SaveNotice(id int64, notice *models.SysNotice) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	notice.CreateBy = user.UserName
	notice.CreateTime = time.Now()
	return mysql.SaveNotice(notice)
}

func UploadNotice(id int64, notice *models.SysNotice) error {
	user, err := mysql.FindUserById(id)
	if err != nil {
		return err
	}
	no, _ := mysql.FindNoticeInfoById(notice.NoticeId)

	if id != 1 && user.UserName != no.CreateBy {
		return errors.New("没有权限进行操作")
	}
	notice.UpdateBy = user.UserName
	notice.UpdateTime = time.Now()
	return mysql.UploadNotice(notice)
}

func DeleteNotice(noticeid int, userid int64) error {
	user, err := mysql.FindUserById(userid)
	if err != nil {
		return err
	}
	no, err := mysql.FindNoticeInfoById(noticeid)
	if err != nil {
		return err
	}
	if userid != 1 && user.UserName != no.CreateBy {
		return errors.New("没有权限进行操作")
	}
	return mysql.DeleteNotice(noticeid)
}
