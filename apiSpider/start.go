package apiSpider

import (
	"github.com/PeterYangs/article-spider/form"
	ff "github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/article-spider/mode"
	"github.com/PeterYangs/article-spider/spider"
)

func Start(form form.Form) {

	spider.SpiderInit(form, mode.Api, func(f ff.Form) {

		GetList(f)
	})

}
