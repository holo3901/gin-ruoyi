package logic

import (
	"ruoyi/dao/mysql"
	"ruoyi/models"
	"time"
)

func SaveContent(content *models.Content) {
	content.CreatedAt = time.Now()

	mysql.SaveContent(content)
}
