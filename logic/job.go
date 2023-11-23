package logic

import (
	"github.com/jinzhu/copier"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SelectJobList(param *models.SearchTableDataParam, job *models.SysJob) ([]*models.SysJob, error) {
	return SelectJobList(param, job)
}

func FindJobById(id int) (*models.SysJob, error) {
	return mysql.FindJobById(id)
}

func SaveJob(param *models.SysJobParam, userid int64) error {
	user, err := mysql.FindUserById(userid)
	if err != nil {
		return err
	}
	job := new(models.SysJob)
	err = copier.Copy(&job, param)
	if err != nil {
		return err
	}
	job.CreateBy = user.UserName
	job.CreateTime = time.Now()
	return mysql.SaveJob(job)

}

func UploadJob(param *models.SysJobParam, userid int64) error {
	user, err := mysql.FindUserById(userid)
	if err != nil {
		return err
	}
	job := new(models.SysJob)
	err = copier.Copy(&job, param)
	if err != nil {
		return err
	}
	job.UpdateBy = user.UserName
	job.UpdateTime = time.Now()
	return mysql.UploadJob(job)
}

func ChangeStatus(id string, status string) error {
	return mysql.ChangeStatus(id, status)
}

func DeleteJob(id string) error {
	return mysql.DeleteJob(id)
}

func SelectJobLogList(param *models.SearchTableDataParam, log *models.SysJobLog) ([]*models.SysJobLog, error) {
	return mysql.SelectJobLogList(param, log)
}

func FindJobLogById(id int) (*models.SysJobLog, error) {
	return mysql.FindJobLogById(id)
}
func DeleteJobLog(id int) error {
	return mysql.DeleteJobLog(id)
}

func ClearJobLog() error {
	return mysql.ClearJobLog()
}
