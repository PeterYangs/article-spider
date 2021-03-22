package connect

import (
	"github.com/gorilla/websocket"
	"sync"
)

//type Connect sync.Map

//存储全局websocket连接
var ConnectList sync.Map

//添加一个连接
func AddCon(uid string, con *websocket.Conn) {

	ConnectList.Store(uid, con)

}

//获取一个连接
func GetCon(uid string) *websocket.Conn {

	con, ok := ConnectList.Load(uid)

	if ok == false {

		return nil
	}

	return con.(*websocket.Conn)

}

//删除一个连接
func DeleteCon(uid string) {

	ConnectList.Delete(uid)

}
