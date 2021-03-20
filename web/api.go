package web

import (
	"article-spider/common"
	"article-spider/form"
	"article-spider/spider"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

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

	r.POST("/submit", func(context *gin.Context) {

		//post:=context.Request.ParseForm
		//
		json := make(map[string]interface{}) //注意该结构接受的内容
		err := context.BindJSON(&json)

		if err != nil {

			fmt.Println(err)

			return
		}

		limit, err := strconv.Atoi(json["limit"].(string))

		if err != nil {

			fmt.Println(err)

			return
		}

		pageStart, err := strconv.Atoi(json["pageStart"].(string))

		if err != nil {

			fmt.Println(err)

			return
		}

		//解析列表选择器和详情选择器
		detailFields := common.ResolveFields((json["detailFields"]).(map[string]interface{}))

		listFields := common.ResolveFields((json["listFields"]).(map[string]interface{}))

		f := form.Form{
			Host:             (json["host"]).(string),
			Channel:          (json["channel"]).(string),
			Limit:            limit,
			PageStart:        pageStart,
			ListSelector:     (json["listSelector"]).(string),
			ListHrefSelector: (json["listHrefSelector"]).(string),
			DetailFields:     detailFields,
			ListFields:       listFields,
			ProxyAddress:     (json["proxyAddress"]).(string),
		}

		go spider.Start(f)

	})

	r.Run(":8089")

}
