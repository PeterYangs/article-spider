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

}
