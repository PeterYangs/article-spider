package spider

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PeterYangs/article-spider/v2/common"
	"github.com/PeterYangs/article-spider/v2/form"
	ff "github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/mode"
	"github.com/PeterYangs/article-spider/v2/result"
	"github.com/PeterYangs/tools"
	uuid "github.com/satori/go.uuid"
	"os"
	"sync"
)

func SpiderInit(form form.Form, mode mode.Mode, listFunc func(form ff.Form)) {

	//设置模式
	form.Mode = mode

	form.Client = common.GetHttpClient(form)

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

	//数据存储管道赋值
	form.Storage = storage

	//日志管道初始化
	form.BroadcastChan = make(chan map[string]string, 3)

	//通知等待锁
	var BroadcastWait sync.WaitGroup

	//通知等待锁赋值
	form.BroadcastWait = &BroadcastWait

	form.BroadcastWait.Add(1)

	//http设置初始化
	//form.HttpSetting = tools.HttpSetting{ProxyAddress: form.ProxyAddress, Header: form.HttpHeader}

	//excel等待锁
	var excelWait sync.WaitGroup

	form.StorageWait = &excelWait

	form.StorageWait.Add(1)

	//进度值初始化
	form.Progress = &sync.Map{}

	//爬取完成通知
	form.IsFinish = make(chan bool, 1)

	if form.ResultCallback != nil {

		go func(f ff.Form) {

			result.GetResult(form, func(item map[string]string) {

				f.ResultCallback(item)
			})

		}(form)

	} else {

		//协程写入Excel
		go WriteExcel(form)

	}

	//协程开启日志输出
	go Broadcast(form)

	//爬取列表
	//GetList(form)
	listFunc(form)

	//等待管道处理完excel写入
	form.StorageWait.Wait()

	if form.ResultCallback == nil {

		uuidString := uuid.NewV4().String()

		//创建excel表文件夹
		tools.MkDirDepth("web/static/excel")

		filename := "web/static/excel/" + uuidString + ".xlsx"

		//发送excel路径
		form.BroadcastChan <- map[string]string{"types": "finish", "data": "static/excel/" + uuidString + ".xlsx"}

		//生成excel
		err = f.SaveAs(filename)

		if err != nil {

			fmt.Println(err)

		}

	}

	//关闭通知管道
	close(form.BroadcastChan)

	//等待通知管道处理完毕
	form.BroadcastWait.Wait()

	fmt.Println("执行完毕")

}
