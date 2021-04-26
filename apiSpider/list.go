package apiSpider

import (
	"github.com/PeterYangs/article-spider/common"
	"github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/article-spider/spider"
	"github.com/PeterYangs/tools"
	"sync"
)

func GetList(form form.Form) {

	//当前页码
	//var pageCurrent int

	hostLast := tools.SubStr(form.Host, len(form.Host)-1, 1)

	if hostLast == "/" {

		form.Host = tools.SubStr(form.Host, 0, len(form.Host)-1)
	}

	ChannelFirst := tools.SubStr(form.Channel, 0, 1)

	if ChannelFirst != "/" {

		form.Channel = "/" + form.Channel
	}

	common.GetChannelList(form, func(listUrl string) {

		//获取html页面
		apiResult, err := tools.GetToString(listUrl, form.HttpSetting)

		detailList := form.ApiConversion(apiResult)

		if err != nil {

			common.ErrorLine(form, err.Error())

			return

		}

		//详情页面并发同步锁
		var wait sync.WaitGroup

		//控制详情页面最大并发数管道
		detailMaxChan := make(chan int, form.DetailMaxCoroutine)

		for _, s := range detailList {

			wait.Add(1)

			//控制最大并发
			if form.DetailMaxCoroutine != 0 {

				detailMaxChan <- 1

			}

			s = common.GetHref(s, form.Host)

			//根据列表的长度开启协程爬取详情页
			go spider.GetDetail(form, s, &wait, detailMaxChan)

		}

		wait.Wait()

	})

	//通知excel已完成
	form.IsFinish <- true

}
