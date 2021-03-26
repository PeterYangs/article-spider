package spider

import (
	"article-spider/connect"
	"article-spider/form"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func broadcast(form form.Form) {

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

				fmt.Println(message["data"])

			case "finish":

				log.Println(message["data"])

			case "error":

				log.Println(message["data"])

			}

		}

		//fmt.Println(message["types"])

	}

}
