package web

import (
	"article-spider/common"
	"article-spider/connect"
	"article-spider/form"
	"article-spider/message"
	"article-spider/spider"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

//启动web服务
func StartWeb() {

	r := gin.Default()

	r.Delims("{[{", "}]}")

	r.Static("/static", "./web/static")

	r.LoadHTMLGlob("./web/html/*")

	r.GET("/ping", func(context *gin.Context) {

		context.JSON(200, gin.H{"msg": "success"})

	})

	r.GET("/", func(context *gin.Context) {

		context.HTML(200, "index.html", gin.H{})
	})

	r.GET("/websocket", func(context *gin.Context) {

		context.HTML(200, "websocket.html", gin.H{})
	})

	//websocket
	r.Any("/broadcast", func(context *gin.Context) {

		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)

		if err != nil {

			fmt.Println(err)

			return
		}

		var uid string

		defer func() {

			err := conn.Close()

			if err != nil {

				fmt.Println(err)
			}

			//删除掉这个连接
			connect.DeleteCon(uid)

		}()

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()

			_ = msgType

			//fmt.Println(msgType)

			if err != nil {

				fmt.Println(err)

				fmt.Println(connect.ConnectList)

				return
			}

			var m message.Message

			err = json.Unmarshal(msg, &m)

			if err != nil {

				fmt.Println(err)

				return
			}

			// Print the message to the console
			//fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			//fmt.Println(m)

			switch m.Types {

			case "registered":

				uid = uuid.NewV4().String()

				connect.AddCon(uid, conn)

				//err:=conn.WriteMessage(msgType, []byte(uid))

				err := conn.WriteJSON(gin.H{"types": m.Types, "data": uid})

				if err != nil {

					fmt.Println(err)

					return
				}

			}

			// Write message back to browser
			//if err = conn.WriteMessage(msgType, msg); err != nil {
			//	return
			//}
		}

	})

	//websocket测试
	r.Any("/echo", func(context *gin.Context) {

		//context.Request.Response
		//conn, err := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(context.Writer, context.Request, nil) // error ignored for sake of simplicity

		if err != nil {

			fmt.Println(err)

			return
		}

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()

			fmt.Println(msgType)

			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}

	})

	r.POST("/submit", func(context *gin.Context) {

		//post:=context.Request.ParseForm
		//
		json_ := make(map[string]interface{}) //注意该结构接受的内容
		err := context.BindJSON(&json_)

		if err != nil {

			fmt.Println(err)

			return
		}

		limit, err := strconv.Atoi(json_["limit"].(string))

		if err != nil {

			fmt.Println(err)

			return
		}

		pageStart, err := strconv.Atoi(json_["pageStart"].(string))

		if err != nil {

			fmt.Println(err)

			return
		}

		//解析列表选择器和详情选择器
		detailFields := common.ResolveFields((json_["detailFields"]).(map[string]interface{}))

		listFields := common.ResolveFields((json_["listFields"]).(map[string]interface{}))

		f := form.Form{
			Host:             (json_["host"]).(string),
			Channel:          (json_["channel"]).(string),
			Limit:            limit,
			PageStart:        pageStart,
			ListSelector:     (json_["listSelector"]).(string),
			ListHrefSelector: (json_["listHrefSelector"]).(string),
			DetailFields:     detailFields,
			ListFields:       listFields,
			ProxyAddress:     (json_["proxyAddress"]).(string),
			Uid:              (json_["uid"]).(string),
		}

		go spider.Start(f)

	})

	r.Run(":8089")

}

//跨域检查
func checkOrigin(r *http.Request) bool {

	return true
}
