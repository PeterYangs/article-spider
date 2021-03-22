package spider

import (
	"article-spider/connect"
	"article-spider/form"
	"fmt"
	"github.com/gin-gonic/gin"
)

func broadcast(form form.Form) {

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

				}

			}

		}

		fmt.Println(message)

	}

}