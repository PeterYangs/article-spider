package spider

import (
	"fmt"
	"github.com/PeterYangs/article-spider/connect"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/tools"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func Broadcast(form form.Form) {

	defer form.BroadcastWait.Done()

	for message := range form.BroadcastChan {

		if form.Uid != "" {

			con := connect.GetCon(form.Uid)

			if con != nil {

				switch message["types"] {

				case "log":

					//发送日志
					con.WriteJSON(gin.H{"types": "log", "data": message["data"]})

				case "finish":

					//已完成
					con.WriteJSON(gin.H{"types": "finish", "data": message["data"]})

				case "error":

				}

			}

		}

		if form.DisableDebug == false {

			switch message["types"] {

			case "log":

				maxPage, _ := form.Progress.Load("maxPage")
				currentPage, _ := form.Progress.Load("currentPage")

				logs := strings.TrimSpace(message["data"])

				logs = strings.Replace(logs, "\n", "", -1)
				logs = strings.Replace(logs, "\r", "", -1)
				logs = strings.Replace(logs, "\r\n", "", -1)

				reset := string([]byte{27, 91, 48, 109})

				red := string([]byte{27, 91, 51, 49, 109})

				fmt.Printf("\r %s, %s当前进度：%d%%%s", tools.SubStr(logs, 0, 60), red, int((currentPage.(float32)/maxPage.(float32))*100), reset)

			case "finish":

				log.Println(message["data"])

			case "error":

				log.Println(message["data"])

			}

		}

		//fmt.Println(message["types"])

	}

}
