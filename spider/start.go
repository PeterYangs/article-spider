package spider

import (
	"github.com/PeterYangs/article-spider/form"
	ff "github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/article-spider/mode"
)

func Start(form form.Form) {

	SpiderInit(form, mode.Normal, func(f ff.Form) {

		GetList(f)
	})

	//start.Start(form,mode.Normal)

	////设置模式
	//form.Mode = mode.Normal
	//
	//form.Client = common.GetHttpClient(form)
	//
	////新建xlsx文件
	//f := excelize.NewFile()
	//
	//// 设置工作簿的默认工作表
	//f.SetActiveSheet(f.NewSheet("Sheet1"))
	//
	////Excel文件对象赋值
	//form.ExcelFile = f
	//
	////数据存储管道初始化
	//storage := make(chan map[string]string, 10)
	//
	////创建图片文件夹
	//err := os.Mkdir("image", 766)
	//
	//if err != nil {
	//
	//	fmt.Println(err)
	//}
	//
	////数据存储管道赋值
	//form.Storage = storage
	//
	////日志管道初始化
	//form.BroadcastChan = make(chan map[string]string, 3)
	//
	////通知等待锁
	//var BroadcastWait sync.WaitGroup
	//
	////通知等待锁赋值
	//form.BroadcastWait = &BroadcastWait
	//
	//form.BroadcastWait.Add(1)
	//
	////http设置初始化
	//form.HttpSetting = tools.HttpSetting{ProxyAddress: form.ProxyAddress, Header: form.HttpHeader}
	//
	////excel等待锁
	//var excelWait sync.WaitGroup
	//
	//form.StorageWait = &excelWait
	//
	//form.StorageWait.Add(1)
	//
	////进度值初始化
	//form.Progress = &sync.Map{}
	//
	////爬取完成通知
	//form.IsFinish = make(chan bool, 1)
	//
	//if form.ResultCallback != nil {
	//
	//	go func(f ff.Form) {
	//
	//		result.GetResult(form, func(item map[string]string) {
	//
	//			f.ResultCallback(item)
	//		})
	//
	//	}(form)
	//
	//} else {
	//
	//	//协程写入Excel
	//	go WriteExcel(form)
	//
	//}
	//
	////协程开启日志输出
	//go Broadcast(form)
	//
	////爬取列表
	//GetList(form)
	//
	////等待管道处理完excel写入
	//form.StorageWait.Wait()
	//
	//uuidString := uuid.NewV4().String()
	//
	////创建excel表文件夹
	//tools.MkDirDepth("web/static/excel")
	//
	//filename := "web/static/excel/" + uuidString + ".xlsx"
	//
	////发送excel路径
	//form.BroadcastChan <- map[string]string{"types": "finish", "data": "static/excel/" + uuidString + ".xlsx"}
	//
	////关闭通知管道
	//close(form.BroadcastChan)
	//
	////等待通知管道处理完毕
	//form.BroadcastWait.Wait()
	//
	////生成excel
	//err = f.SaveAs(filename)
	//
	//if err != nil {
	//
	//	fmt.Println(err)
	//
	//}
	//
	//fmt.Println("执行完毕")

}
