package chromedpSpider

import (
	"github.com/PeterYangs/article-spider/v2/form"
	ff "github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/mode"
	"github.com/PeterYangs/article-spider/v2/spider"
)

func Start(form form.Form) {

	spider.SpiderInit(form, mode.Auto, func(f ff.Form) {

		GetList(f)
	})

}
