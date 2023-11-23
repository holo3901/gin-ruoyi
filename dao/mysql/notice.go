package mysql

import (
	"fmt"
	"ruoyi/models"
)

func SelectSysNoticeList(param *models.SearchTableDataParam, notice *models.SysNotice) ([]*models.SysNotice, error) {
	sqlstr := `select * from sys_notice where 1=1`
	var args []interface{}

	if notice.NoticeTitle != "" {
		sqlstr += ` AND notice_title like ?`
		args = append(args, "%"+notice.NoticeTitle+"%")
	}
	if notice.NoticeType != "" {
		sqlstr += ` AND notice_type =?`
		args = append(args, notice.NoticeType)
	}
	if notice.CreateBy != "" {
		sqlstr += ` AND create_by like ?`
		args = append(args, "%"+notice.CreateBy+"%")
	}

	if param.Params.BeginTime != "" {
		start, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND create_time>=? AND create_time <=?`
		args = append(args, start, end)
	}
	sqlstr += ` order by notice_id`
	if param.PageSize != 0 {
		sqlstr += ` LIMIT ? OFFSET ?`
		args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	}
	list := make([]*models.SysNotice, 0)
	fmt.Println(sqlstr, args)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindNoticeInfoById(id int) (*models.SysNotice, error) {
	sqlstr := `select * from sys_notice where notice_id =?`
	a := new(models.SysNotice)
	err := db.Get(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func SaveNotice(notice *models.SysNotice) error {
	sqlstr := `insert into sys_notice (notice_title,notice_type,notice_content,create_by,create_time,remark) 
               values(?,?,?,?,?,?)`
	_, err := db.Exec(sqlstr, notice.NoticeTitle, notice.NoticeType, notice.NoticeContent, notice.CreateBy, notice.CreateTime, notice.Remark)
	if err != nil {
		return err
	}
	return nil
}
func UploadNotice(notice *models.SysNotice) error {
	sqlstr := `update sys_notice set notice_title=?,notice_type=?,notice_content=?,update_by=?,update_time=?,remark=?where notice_id =?`
	_, err := db.Exec(sqlstr, notice.NoticeTitle, notice.NoticeType, notice.NoticeContent, notice.UpdateBy, notice.UpdateTime, notice.Remark, notice.NoticeId)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNotice(id int) error {
	sqlstr := `delete from sys_notice where notice_id =?`
	_, err := db.Exec(sqlstr, id)
	if err != nil {
		return err
	}
	return nil
}
