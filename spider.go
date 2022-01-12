package article_spider

import (
	"context"
	"fmt"
	"github.com/PeterYangs/request"
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
}

func NewSpider(f Form, mode Mode) *Spider {

	client := request.NewClient()

	if f.HttpTimeout != 0 {

		client.Timeout(f.HttpTimeout)
	}

	client.Header(f.HttpHeader)

	client.ReTry(1)

	if f.HttpProxy != "" {

		client.Proxy(f.HttpProxy)

	}

	detailMaxCoroutines := 30

	//如果手动设置的详情协程数大于最大详情协程数或者等于0，则将设置成最大协程数
	if f.DetailCoroutineNumber < detailMaxCoroutines && f.DetailCoroutineNumber != 0 {

		detailMaxCoroutines = f.DetailCoroutineNumber

	}

	cxt, cancel := context.WithCancel(context.Background())

	return &Spider{form: f, mode: mode, client: client, detailCoroutineChan: make(chan bool, detailMaxCoroutines), cxt: cxt, cancel: cancel, wait: sync.WaitGroup{}, imageDir: "image", detailWait: sync.WaitGroup{}}
}

func (s *Spider) Debug() *Spider {

	s.debug = true

	return s
}

func (s *Spider) Start() {

	s.notice = NewNotice(s)

	s.result = NewResult(s)

	s.excel = NewExcel(s)

	go s.notice.Service()

	go s.result.Work()

	s.form.s = s

	switch s.mode {

	case Normal:

		NewNormal(s).Start()

	}

	s.wait.Wait()

	fmt.Println("finish")

	//select {}

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

			callback(i)

		}

	}

}
