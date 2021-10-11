package spider

import (
	"github.com/PeterYangs/article-spider/v2/conf"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/mode"
	"github.com/PeterYangs/article-spider/v2/mode/api"
	"github.com/PeterYangs/article-spider/v2/mode/normal"
	"github.com/PeterYangs/article-spider/v2/notice"
	"github.com/PeterYangs/article-spider/v2/result"
	"github.com/PeterYangs/request"
	"github.com/PeterYangs/tools"
	"strconv"
	"strings"
	"sync"
)

type Spider struct {
	form   *form.Form
	Notice *notice.Notice
}

func NewSpider() *Spider {

	return &Spider{Notice: notice.NewNotice()}
}

func (s *Spider) LoadForm(cf form.CustomForm) *Spider {

	f := &form.Form{
		Host:                       cf.Host,
		Channel:                    cf.Channel,
		PageStart:                  cf.PageStart,
		Length:                     cf.Length,
		ListSelector:               cf.ListSelector,
		HrefSelector:               cf.HrefSelector,
		DisableAutoCoding:          cf.DisableAutoCoding,
		LazyImageAttrName:          cf.LazyImageAttrName,
		DisableImageExtensionCheck: cf.DisableImageExtensionCheck,
		AllowImageExtension:        cf.AllowImageExtension,
		DefaultImg:                 cf.DefaultImg,
		DetailFields:               cf.DetailFields,
		ListFields:                 cf.ListFields,
		CustomExcelHeader:          cf.CustomExcelHeader,
		DetailCoroutineNumber:      cf.DetailCoroutineNumber,
		HttpTimeout:                cf.HttpTimeout,
		HttpHeader:                 cf.HttpHeader,
		MiddleSelector:             cf.MiddleHrefSelector,
		ResultCallback:             cf.ResultCallback,
		ApiConversion:              cf.ApiConversion,
	}

	s.form = f

	//通知服务
	s.form.Notice = s.Notice

	//初始化结果通道
	s.form.Storage = make(chan map[string]string, 20)

	s.form.Wait = sync.WaitGroup{}

	return s
}

func (s *Spider) loadClient() *Spider {

	client := request.NewClient()

	if s.form.HttpTimeout != 0 {

		client.Timeout(s.form.HttpTimeout)
	}

	client.Header(s.form.HttpHeader)

	client.ReTry(1)

	s.form.Client = client

	return s

}

func (s *Spider) StartApi() {

	s.form.Mode = mode.Api

	detailMaxCoroutines := conf.Conf.DetailMaxCoroutines

	//如果手动设置的详情协程数大于最大详情协程数或者等于0，则将设置成最大协程数
	if s.form.DetailCoroutineNumber > detailMaxCoroutines || s.form.DetailCoroutineNumber == 0 {

		s.form.DetailCoroutineNumber = detailMaxCoroutines
	}

	s.form.DetailCoroutineChan = make(chan bool, s.form.DetailCoroutineNumber)

	//消息处理服务
	go s.Notice.Service(func() {

		s.form.Wait.Done()
	})

	r := result.NewResult(s.form)

	//excel处理等待标记
	s.form.Wait.Add(1)

	//处理结果服务
	go r.Work()

	s.checkLink()

	//初始化http客户端
	s.loadClient()

	//消息关闭等待标记
	s.form.Wait.Add(1)

	n := api.NewApi(s.form)

	//列表回调
	s.getChannelList(func(listUrl string) {

		n.GetList(listUrl)

	})

	//等待详情协程处理完毕
	s.form.DetailWait.Wait()

	close(s.form.Storage)

	s.form.Wait.Wait()

}

// Start 普通模式爬取
func (s *Spider) Start() {

	s.form.Mode = mode.Normal

	detailMaxCoroutines := conf.Conf.DetailMaxCoroutines

	//如果手动设置的详情协程数大于最大详情协程数或者等于0，则将设置成最大协程数
	if s.form.DetailCoroutineNumber > detailMaxCoroutines || s.form.DetailCoroutineNumber == 0 {

		s.form.DetailCoroutineNumber = detailMaxCoroutines
	}

	s.form.DetailCoroutineChan = make(chan bool, s.form.DetailCoroutineNumber)

	//消息处理服务
	go s.Notice.Service(func() {

		s.form.Wait.Done()
	})

	r := result.NewResult(s.form)

	//excel处理等待标记
	s.form.Wait.Add(1)

	//处理结果服务
	go r.Work()

	s.checkLink()

	//初始化http客户端
	s.loadClient()

	//消息关闭等待标记
	s.form.Wait.Add(1)

	n := normal.NewNormal(s.form)

	//列表回调
	s.getChannelList(func(listUrl string) {

		n.GetList(listUrl)

	})

	//等待详情协程处理完毕
	s.form.DetailWait.Wait()

	close(s.form.Storage)

	s.form.Wait.Wait()
}

func (s *Spider) checkLink() {

	hostLast := tools.SubStr(s.form.Host, len(s.form.Host)-1, 1)

	if hostLast == "/" {

		s.form.Host = tools.SubStr(s.form.Host, 0, len(s.form.Host)-1)
	}

	ChannelFirst := tools.SubStr(s.form.Channel, 0, 1)

	if ChannelFirst != "/" {

		s.form.Channel = "/" + s.form.Channel
	}

}

// GetChannelList 获取栏目链接
func (s *Spider) getChannelList(callback func(listUrl string)) {

	switch s.form.Mode {

	case mode.Normal, mode.Api:

		//if form.ChannelFunc == nil {

		//当前页码
		var pageCurrent int

		//form.Progress.Store("maxPage", float32(form.PageStart+form.Limit))
		//form.Progress.Store("currentPage", float32(0))

		for pageCurrent = s.form.PageStart; pageCurrent < s.form.PageStart+s.form.Length; pageCurrent++ {

			//当前列表url
			url := s.form.Host + strings.Replace(s.form.Channel, "[PAGE]", strconv.Itoa(pageCurrent), -1)

			callback(url)

			//currentPage, _ := form.Progress.Load("currentPage")

			////这里有点恶心，有没有简单的写法
			//c := currentPage.(float32)
			//c++
			//form.Progress.Store("currentPage", c)

		}

		return
		//}

		////自定义栏目
		//for _, i := range form.ChannelFunc(form) {
		//
		//	callback(form.Host + i)
		//
		//	currentPage, _ := form.Progress.Load("currentPage")
		//
		//	c := currentPage.(float32)
		//	c++
		//	form.Progress.Store("currentPage", c)
		//
		//}

		//case mode.Auto:
		//
		//	//当前页码
		//	var pageCurrent int
		//
		//	form.Progress.Store("maxPage", float32(form.Limit))
		//	form.Progress.Store("currentPage", float32(0))
		//
		//	for pageCurrent = 0; pageCurrent < form.Limit; pageCurrent++ {
		//
		//		callback(strconv.Itoa(pageCurrent))
		//
		//		currentPage, _ := form.Progress.Load("currentPage")
		//
		//		c := currentPage.(float32)
		//		c++
		//		form.Progress.Store("currentPage", c)
		//
		//	}

	}

}
