package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jianfengye/collection"
	"go.uber.org/zap"
	"net/http"
	"ruoyi/logic"
	"ruoyi/models"
	"ruoyi/pkg/safe"
	"strconv"
	"sync"
	"time"
)

// 客户端连接详情
type wsClients struct {
	Conn *websocket.Conn `json:"conn"`

	RemoteAddr string `json:"remote_addr"`

	Uid string `json:"uid"`

	Username string `json:"username"`

	RoomId string `json:"room_id"`

	AvatarId string `json:"avatar_id"`
}

type msgData struct {
	Uid      string        `json:"uid"`
	Username string        `json:"username"`
	AvatarId string        `json:"avatar_id"`
	ToUid    string        `json:"to_uid"`
	Content  string        `json:"content"`
	ImageUrl string        `json:"image_url"`
	RoomId   string        `json:"room_id"`
	Count    int           `json:"count"`
	List     []interface{} `json:"list"`
	Time     int64         `json:"time"`
}

// client & serve 的消息体
type msg struct {
	Status int             `json:"status"`
	Data   msgData         `json:"data"`
	Conn   *websocket.Conn `json:"conn"`
}

type pingStorage struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Time       int64           `json:"time"`
}

// 变量定义初始化
var (
	wsUpgrader = websocket.Upgrader{}

	clientMsg = msg{}

	mutex = sync.Mutex{}

	//rooms = [roomCount + 1][]wsClients{}
	rooms = make(map[int][]interface{})

	enterRooms = make(chan wsClients)

	sMsg = make(chan msg)

	offline = make(chan *websocket.Conn)

	chNotify = make(chan int, 1)

	pingMap []interface{}
)

// 定义消息类型
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTypeSend = 3          // 消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 6 // 房间总数

func Run(ctx *gin.Context) {
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, _ := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	defer c.Close()
	done := make(chan struct{})

	go read(c, done)
	go write(done)

	select {}
}

func read(c *websocket.Conn, done chan<- struct{}) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("read 发生错误", zap.Error(errors.New("")))
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			offline <- c
			zap.L().Error("read 发生错误", zap.Error(err))
			c.Close()
			close(done)
			return
		}

		serveMsgStr := message

		// 心跳处理响应，heartbeat为与客户端约定的值
		if string(serveMsgStr) == `heartbeat` {
			appendPing(c)
			chNotify <- 1
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":0,"data":"heartbeat ok"}`))
			<-chNotify
			continue
		}

		json.Unmarshal(message, &clientMsg)

		if clientMsg.Data.Uid != "" {
			if clientMsg.Status == msgTypeOnline { //进去房间，建立连接
				roomId, _ := getRoomId()

				enterRooms <- wsClients{
					Conn:       c,
					RemoteAddr: c.RemoteAddr().String(),
					Uid:        clientMsg.Data.Uid,
					Username:   clientMsg.Data.Username,
					RoomId:     roomId,
					AvatarId:   clientMsg.Data.AvatarId,
				}
			}

			_, serveMsg := formatServeMsgStr(clientMsg.Status, c)
			sMsg <- serveMsg
		}
	}
}

func write(done <-chan struct{}) {

	defer func() {
		if err := recover(); err != nil {
			zap.L().Error("write 发生错误", zap.Error(errors.New("")))
		}
	}()

	for {
		select {
		case <-done:
			return
		case r := <-enterRooms:
			handleConnClients(r.Conn)
		case cl := <-sMsg:
			serverMsgStr, _ := json.Marshal(cl)
			switch cl.Status {
			case msgTypeOnline, msgTypeSend:
				notify(cl.Conn, string(serverMsgStr))
			case msgTypeGetOnlineUser:
				chNotify <- 1
				cl.Conn.WriteMessage(websocket.TextMessage, serverMsgStr)
				<-chNotify
			case msgTypePrivateChat:
				chNotify <- 1
				toc := findToUserCoonClient()
				if toc != nil {
					toc.(wsClients).Conn.WriteMessage(websocket.TextMessage, serverMsgStr)
				}
				<-chNotify
			}
		case o := <-offline:
			disconnect(o)
		}
	}
}

func handleConnClients(c *websocket.Conn) {
	roomId, roomIdInt := getRoomId()

	objColl := collection.NewObjCollection(rooms[roomIdInt])

	retColl := safe.Safety.Do(func() interface{} {
		return objColl.Reject(func(item interface{}, key int) bool {
			if item.(wsClients).Uid == clientMsg.Data.Uid {
				chNotify <- 1
				item.(wsClients).Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
				<-chNotify
				return true
			}
			return false
		})
	}).(collection.ICollection)

	retColl = safe.Safety.Do(func() interface{} {
		return retColl.Append(wsClients{
			Conn:       c,
			RemoteAddr: c.RemoteAddr().String(),
			Uid:        clientMsg.Data.Uid,
			Username:   clientMsg.Data.Username,
			RoomId:     roomId,
			AvatarId:   clientMsg.Data.AvatarId,
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()

	rooms[roomIdInt] = interfaces
}

// 统一消息发放
func notify(conn *websocket.Conn, msg string) {
	chNotify <- 1
	_, roomIdInt := getRoomId()
	assignRoom := rooms[roomIdInt]
	for _, con := range assignRoom {
		if con.(wsClients).RemoteAddr != conn.RemoteAddr().String() {
			con.(wsClients).Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
	<-chNotify
}
func appendPing(c *websocket.Conn) {
	objColl := collection.NewObjCollection(pingMap)

	//先删除相同的
	retColl := safe.Safety.Do(func() interface{} {
		return objColl.Reject(func(obj interface{}, index int) bool {
			if obj.(pingStorage).RemoteAddr == c.RemoteAddr().String() {
				return true
			}
			return false
		})
	}).(collection.ICollection)

	//再追加
	retColl = safe.Safety.Do(func() interface{} {
		return retColl.Append(pingStorage{
			Conn:       c,
			RemoteAddr: c.RemoteAddr().String(),
			Time:       time.Now().Unix(),
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()
	pingMap = interfaces
}

// 获取私聊的用户连接
func findToUserCoonClient() interface{} {
	_, roomIdInt := getRoomId()
	toUserUid := clientMsg.Data.ToUid
	assignRoom := rooms[roomIdInt]
	for _, c := range assignRoom {
		stringUid := c.(wsClients).Uid
		if stringUid == toUserUid {
			return c
		}
	}
	return nil
}

// 格式化传送给客户端的消息数据
func formatServeMsgStr(status int, conn *websocket.Conn) ([]byte, msg) {
	roomId, roomIdInt := getRoomId()

	data := msgData{
		Username: clientMsg.Data.Username,
		Uid:      clientMsg.Data.Uid,
		RoomId:   roomId,
		Time:     time.Now().UnixNano() / 1e6,
	}

	if status == msgTypeSend || status == msgTypePrivateChat {
		data.AvatarId = clientMsg.Data.AvatarId
		content := clientMsg.Data.Content

		data.Content = content
		if safe.MbStrLen(content) > 800 {
			//直接截断
			data.Content = string([]rune(content)[:800])
		}

		toUidStr := clientMsg.Data.ToUid
		toUid, _ := strconv.Atoi(toUidStr)

		//保存信息
		stringUid := data.Uid
		intUid, _ := strconv.Atoi(stringUid)

		p := new(models.Content)
		p.UserId = intUid
		p.ToUserId = toUid
		p.Content = data.Content
		p.RoomId = data.RoomId
		p.ImageUrl = clientMsg.Data.ImageUrl
		logic.SaveContent(p)
	}

	if status == msgTypeGetOnlineUser {
		ro := rooms[roomIdInt]
		data.Count = len(ro)
		data.List = ro
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
		Conn:   conn,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr, jsonStrServeMsg
}

// 离线通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()
	objColl := collection.NewObjCollection(rooms[roomIdInt])

	retColl := safe.Safety.Do(func() interface{} {
		return objColl.Reject(func(item interface{}, key int) bool {
			if item.(wsClients).RemoteAddr == conn.RemoteAddr().String() {

				data := msgData{
					Username: item.(wsClients).Username,
					Uid:      item.(wsClients).Uid,
					Time:     time.Now().UnixNano() / 1e6,
				}

				jsonStrServeMsg := msg{
					Status: msgTypeOffline,
					Data:   data,
				}
				serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

				disMsg := string(serveMsgStr)

				item.(wsClients).Conn.Close()

				notify(conn, disMsg)

				return true
			}
			return false
		})
	}).(collection.ICollection)

	interfaces, _ := retColl.ToInterfaces()
	rooms[roomIdInt] = interfaces
}

func getRoomId() (string, int) {
	roomId := clientMsg.Data.RoomId

	roomIdInt, _ := strconv.Atoi(roomId)
	return roomId, roomIdInt
}

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
