package sendemail

import (
	mailer "gopkg.in/gomail.v2"
	"ruoyi/models"
	"ruoyi/settings"
	"strconv"
)

func SendEmail(p *models.SendEmail, rand int) error {
	msg := mailer.NewMessage()
	msg.SetHeader("From", settings.Conf.EmailConfig.SmtpEmail)
	msg.SetHeader("To", msg.FormatAddress(p.Email, p.User))
	msg.SetHeader("Subject", "邮箱验证")
	msg.SetBody("text/html", "【HL】 验证码: "+strconv.Itoa(rand)+",你正在绑定邮箱信息(若非本人操作，请删除本短信)")
	d := mailer.NewDialer(settings.Conf.EmailConfig.SmtpHost, 465, settings.Conf.EmailConfig.SmtpEmail, settings.Conf.EmailConfig.SmtpPass)
	err := d.DialAndSend(msg)
	return err
}
