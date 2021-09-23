package spider

import (
	"github.com/PeterYangs/article-spider/v2/form"
	ff "github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/mode"
)

func Start(form form.Form) {

	SpiderInit(form, mode.Normal, func(f ff.Form) {

		GetList(f)
	})

}
