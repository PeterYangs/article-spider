package web

import "github.com/gin-gonic/gin"

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

	r.Run(":8089")

}
