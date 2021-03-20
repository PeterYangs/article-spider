package spider

import (
	"article-spider/form"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PeterYangs/tools"
	"github.com/satori/go.uuid"
	"os"
	"sync"
)

func Start(form form.Form) {

	//新建xlsx文件
	f := excelize.NewFile()

	// 设置工作簿的默认工作表
	f.SetActiveSheet(f.NewSheet("Sheet1"))

	//Excel文件对象赋值
	form.ExcelFile = f

	//数据存储管道初始化
	storage := make(chan map[string]string, 10)

	//创建图片文件夹
	err := os.Mkdir("image", 766)

	if err != nil {

		fmt.Println(err)
	}

	//管道赋值
	form.Storage = storage

	//http设置初始化
	form.HttpSetting = tools.HttpSetting{ProxyAddress: form.ProxyAddress, Header: form.HttpHeader}

	//excel等待锁
	var excelWait sync.WaitGroup

	form.ExcelWait = &excelWait

	form.ExcelWait.Add(1)

	form.IsFinish = make(chan bool, 1)

	//协程写入Excel
	go WriteExcel(form)

	//爬取列表
	GetList(form)

	//close(form.Storage)

	//等待管道处理完excel写入
	form.ExcelWait.Wait()

	uuidString := uuid.NewV4().String()

	err = f.SaveAs(uuidString + ".xlsx")

	if err != nil {

		fmt.Println(err)

	}

	fmt.Println("执行完毕")

}
