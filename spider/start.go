package spider

import (
	"article-spider/form"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/satori/go.uuid"
)

func Start(form form.Form) {

	//新建xlsx文件
	f := excelize.NewFile()

	// 设置工作簿的默认工作表
	f.SetActiveSheet(f.NewSheet("Sheet1"))

	//赋值
	form.ExcelFile = f

	//管道初始化
	storage := make(chan map[string]string, 10)

	//管道赋值
	form.Storage = storage

	//协程写入Excel
	go WriteExcel(form)

	//爬取列表
	GetList(form)

	close(form.Storage)

	uuidString := uuid.NewV4().String()

	f.SaveAs(uuidString + ".xlsx")

}
