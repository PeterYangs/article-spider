package spider

import (
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/mode"
	"github.com/PeterYangs/article-spider/v2/mode/normal"
	"github.com/PeterYangs/article-spider/v2/notice"
	"github.com/PeterYangs/tools"
	"github.com/PeterYangs/tools/http"
	"strconv"
	"strings"
)

type Spider struct {
	form   *form.Form
	notice *notice.Notice
}

func NewSpider() *Spider {

	return &Spider{notice: notice.NewNotice()}
}

func (s *Spider) LoadForm(form *form.Form) *Spider {

	s.form = form

	return s
}

func (s *Spider) loadClient() *Spider {

	client := http.Client()

	s.form.Client = client

	return s

}

// Start 普通模式爬取
func (s *Spider) Start() {

	s.checkLink()

	s.loadClient()

	n := normal.NewNormal(s.form)

	s.getChannelList(func(listUrl string) {

		n.GetList(listUrl)

	})

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
