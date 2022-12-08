package article_spider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/PeterYangs/request/v2"
	"github.com/PeterYangs/tools"
	"strconv"
	"strings"
	"sync"
)

type Spider struct {
	form                Form
	mode                Mode
	client              *request.Client
	detailCoroutineChan chan bool
	cxt                 context.Context
	cancel              context.CancelFunc
	wait                sync.WaitGroup
	notice              *Notice
	imageDir            string
	result              *result
	excel               *excel
	currentIndex        int
	detailWait          sync.WaitGroup //详情等待
	detailSize          int
	total               int
	autoPage            int //自动化模式当前页码
	debug               bool
	savePath            string                                                         //图片保存文件夹，不会出现在图片路径中，为空则为当前运行路径
	CustomDownloadFun   func(imgUrl string, imgPath string, f *Form, item Field) error //自实现图片下载
}

func NewSpider(f Form, mode Mode, cxt context.Context) *Spider {

	client := request.NewClient()

	if f.HttpTimeout != 0 {

		client.Timeout(f.HttpTimeout)
	}

	tr := client.GetTransport()

	//跳过证书检查
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client.Header(f.HttpHeader)

	client.ReTry(1)

	if f.HttpProxy != "" {

		client.Proxy(f.HttpProxy)

	}

	detailMaxCoroutines := 1

	//设置最大协程数，默认是1
	if f.DetailCoroutineNumber != 0 {

		detailMaxCoroutines = f.DetailCoroutineNumber
	}

	cxt2, cancel := context.WithCancel(cxt)

	return &Spider{form: f, mode: mode, client: client, detailCoroutineChan: make(chan bool, detailMaxCoroutines), cxt: cxt2, cancel: cancel, wait: sync.WaitGroup{}, imageDir: "image", detailWait: sync.WaitGroup{}}
}

// SetImageDir 设置图片文件夹
func (s *Spider) SetImageDir(path string) {

	s.imageDir = path

}

// SetSavePath 图片保存文件夹，不会出现在图片路径中，为空则为当前运行路径
func (s *Spider) SetSavePath(path string) {

	s.savePath = path
}

func (s *Spider) Debug() *Spider {

	s.debug = true

	return s
}

// CustomDownloadImage 自实现图片下载
func (s *Spider) CustomDownloadImage(fun func(imgUrl string, imgPath string, f *Form, item Field) error) *Spider {

	s.CustomDownloadFun = fun

	return s

}

func (s *Spider) Start() error {

	//对输入的域名和栏目进行处理
	hostLast := tools.SubStr(s.form.Host, len(s.form.Host)-1, 1)

	if hostLast == "/" {

		s.form.Host = tools.SubStr(s.form.Host, 0, len(s.form.Host)-1)
	}

	ChannelFirst := tools.SubStr(s.form.Channel, 0, 1)

	if ChannelFirst != "/" {

		s.form.Channel = "/" + s.form.Channel
	}

	s.notice = NewNotice(s)

	defer s.notice.Stop()

	s.result = NewResult(s)

	s.excel = NewExcel(s)

	go s.notice.Service()

	go s.result.Work()

	s.form.s = s

	var err error

	switch s.mode {

	case Normal:

		NewNormal(s).Start()

	case Api:

		NewApi(s).Start()

	case Auto:

		err = NewAuto(s).Start()

	case Url:

		NewUrl(s).Start()

	}

	s.wait.Wait()

	fmt.Println("finish")

	return err

}

//--------------------------------------------------------------------------------------------

// GetChannelList 获取栏目链接
func (s *Spider) getChannelList(callback func(listUrl string)) {

	switch s.mode {

	case Normal, Api:

		if s.form.ChannelFunc == nil {

			//当前页码
			var pageCurrent int

			for pageCurrent = s.form.PageStart; pageCurrent < s.form.PageStart+s.form.Length; pageCurrent++ {

				//当前列表url
				url := s.form.Host + strings.Replace(s.form.Channel, "[PAGE]", strconv.Itoa(pageCurrent), -1)

				callback(url)

			}

			return
		}

		cList := s.form.ChannelFunc(&s.form)

		s.form.Length = len(cList)

		//自定义栏目
		for _, i := range cList {

			ChannelFirst := tools.SubStr(i, 0, 1)

			if ChannelFirst != "/" {

				i = "/" + i
			}

			callback(s.form.Host + i)

		}

	}

}
