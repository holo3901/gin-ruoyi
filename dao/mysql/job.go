package mysql

import (
	"ruoyi/models"
	"strconv"
)

func SelectJobList(param *models.SearchTableDataParam, job *models.SysJob) ([]*models.SysJob, error) {
	sqlstr := `select * from sys_job where 1=1`
	var args []interface{}
	if job.JobName != "" {
		sqlstr += ` AND job_name like ?`
		args = append(args, "%"+job.JobName+"%")
	}
	if job.JobGroup != "" {
		sqlstr += ` AND job_group =?`
		args = append(args, job.JobGroup)
	}
	if job.Status != "" {
		sqlstr += ` AND status =?`
		args = append(args, job.JobGroup)
	}
	if job.InvokeTarget != "" {
		sqlstr += ` AND invoke_target like ?`
		args = append(args, "%"+job.InvokeTarget+"%")
	}
	sqlstr += " order by job_id"
	if param.PageSize != 0 && param.PageNum != 0 {
		sqlstr += " DESC LIMIT ? OFFSET ? "
		args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	}
	list := make([]*models.SysJob, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindJobById(id int) (*models.SysJob, error) {
	sqlstr := `select * from sys_job where job_id =?`
	job := new(models.SysJob)
	err := db.Get(job, sqlstr, strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	return job, nil
}

func SaveJob(job *models.SysJob) error {
	sqlstr := `insert into sys_job (job_name,job_group,invoke_target,cron_expression,misfire_policy,concurrent,status,create_by,create_time,remark)
               values (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := db.Exec(sqlstr, job.JobName, job.JobGroup, job.InvokeTarget, job.CronExpression, job.MisfirePolicy, job.Concurrent, job.Status, job.CreateBy, job.CreateTime, job.Remark)
	return err
}

func UploadJob(job *models.SysJob) error {
	sqlstr := `update sys_job set job_name=?,job_group=?,invoke_target=?,cron_expression=?,misfire_policy=?,concurrent=?,status=?,update_by=?,update_time=?,remark=? where job_id =?  `
	_, err := db.Exec(sqlstr, job.JobName, job.JobGroup, job.InvokeTarget, job.CronExpression, job.MisfirePolicy, job.Concurrent, job.Status, job.UpdateBy, job.UpdateTime, job.Remark, job.JobId)
	return err
}

func ChangeStatus(id string, status string) error {
	sqlstr := `update sys_job set status=? where job_id=?`
	_, err := db.Exec(sqlstr, status, id)
	return err
}

func DeleteJob(id string) error {
	sqlstr := `delete from sys_job where job_id in (?)`
	_, err := db.Exec(sqlstr, id)
	return err
}

func SelectJobLogList(param *models.SearchTableDataParam, log *models.SysJobLog) ([]*models.SysJobLog, error) {
	sqlstr := ` select * from sys_job_log where 1=1 `
	var args []interface{}
	if log.JobName != "" {
		sqlstr += ` AND job_name like ?`
		args = append(args, "%"+log.JobName+"%")
	}
	if log.JobGroup != "" {
		sqlstr += ` AND job_group = ?`
		args = append(args, log.JobGroup)
	}
	if log.Status != "" {
		sqlstr += ` AND status=?`
		args = append(args, log.Status)
	}
	if param.Params.BeginTime != "" {
		satrt, end := models.GetBeginAndEndTime(param.Params.BeginTime, param.Params.EndTime)
		sqlstr += ` AND create_time >= ? AND create_time<=?`
		args = append(args, satrt, end)
	}
	sqlstr += ` order by job_log_id DESC  LIMIT ? OFFSET ?`
	args = append(args, param.PageSize, (param.PageNum-1)*param.PageSize)
	list := make([]*models.SysJobLog, 0)
	err := db.Select(&list, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindJobLogById(id int) (*models.SysJobLog, error) {
	sqlstr := `select * from sys_job_log where job_log_id=?`
	a := new(models.SysJobLog)
	err := db.Get(a, sqlstr, id)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func DeleteJobLog(id int) error {
	sqlstr := `delete from sys_job_log where job_log_id =?`
	_, err := db.Exec(sqlstr, id)
	return err

}

func ClearJobLog() error {
	sqlstr := `TRUNCATE TABLE sys_job_log;`
	_, err := db.Exec(sqlstr)
	return err
}
