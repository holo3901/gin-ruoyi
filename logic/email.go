package logic

import (
	"math/rand"
	"ruoyi/dao/mysql"
	"ruoyi/dao/redis"
	"ruoyi/models"
	"ruoyi/pkg/sendemail"
	"strconv"
)

func SendEmail(email *models.SendEmail, id int64) error {
	randDomNum := rand.Intn(9000) + 1000
	err := sendemail.SendEmail(email, randDomNum)
	if err != nil {
		return err
	}
	err = redis.CreateEmail(id, randDomNum, email.Email)
	if err != nil {
		return err
	}
	return nil
}

func YanzhengEmail(rand int, id int64) (err error) {
	yanzengma, email, err := redis.GetEmail(id)
	if err != nil {
		return err
	}
	if strconv.Itoa(rand) == yanzengma {
		user, _ := mysql.FindUserById(id)
		user.Email = "'" + email + "'"
		mysql.UpdateUser(user)

	}
	return
}
