package mysql

import "ruoyi/models"

func SaveContent(content *models.Content) {
	sqlstr := `insert into message (user_id,room_id,to_user_id,content,image_url,create_at)
               values (?,?,?,?,?,?)`
	db.Exec(sqlstr, content.UserId, content.RoomId, content.ToUserId, content.Content, content.CreatedAt)
}
