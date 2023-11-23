package scheduler

import (
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"ruoyi/scheduler/tesk"
	"strconv"
	"strings"
)

var c *cron.Cron

// 初始化 定时
func InitCron() {
	c = cron.New()
	c.Start()
	openMethod()
	RunSqlCron()
}

func AddCronFunc(sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}

func RunCronFunc(sepc string, invokeTarget string) {
	a := strings.Split(invokeTarget, "(")
	if len(a) > 1 {
		var x []string

		x = strings.Split(strings.Split(a[1], ")")[0], ",")
		AddCronFunc(sepc, func() {
			Call(m, a[0], convertStringSliceToInterfaceSlice(x)...)
		})
	} else {
		AddCronFunc(sepc, func() {
			Call(m, a[0])
		})
	}
}

func convertStringSliceToInterfaceSlice(strSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		if v == "true" || v == "false" {
			x, _ := strconv.ParseBool(v)
			interfaceSlice[i] = x
		} else {
			interfaceSlice[i] = v
		}
	}
	return interfaceSlice
}

var m map[string]interface{}

func openMethod() {
	m = make(map[string]interface{})
	m["ryTask.ryNoParams"] = tesk.NoParamsMethod
	m["ryTask.ryParams"] = tesk.ParamsMethod
	m["ryTask.ryMultipleParams"] = tesk.MultipleParamsMethod
}

func RunSqlCron() {
	p := new(models.SearchTableDataParam)
	a := new(models.SysJob)
	taskList, _ := mysql.SelectJobList(p, a)
	if len(taskList) == 0 {
		return
	}
	for _, item := range taskList {
		var policy = item.MisfirePolicy

		concurrent := item.Concurrent
		invokeTarget := item.InvokeTarget
		expression := item.CronExpression
		// 获取参数
		if concurrent == 0 {
			if policy == 1 {
				fmt.Println(expression, invokeTarget)
				RunCronFunc(expression, invokeTarget)
			} else if policy == 2 {
				Call(m, invokeTarget)
			}
		}
	}
}

func RunOne(job *models.SysJob) {
	if job.Concurrent == 0 {
		if job.MisfirePolicy == 1 {
			RunCronFunc(job.CronExpression, job.InvokeTarget)
		} else if job.MisfirePolicy == 2 {
			Call(m, job.InvokeTarget)
		}
	}
}
