package result

import (
	"fmt"
	"github.com/PeterYangs/article-spider/v2/conf"
	"github.com/PeterYangs/article-spider/v2/excel"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/notice"
	"github.com/PeterYangs/tools"
	"github.com/shopspring/decimal"
)

type result struct {
	form *form.Form
}

func NewResult(form *form.Form) *result {

	return &result{form: form}
}

func (r *result) Work() {

	defer func() {

		r.form.Wait.Done()

		r.form.Notice.Close()

	}()

	exc := excel.NewExcel(r.form)

	for s := range r.form.Storage {

		if r.form.ResultCallback != nil {

			r.form.ResultCallback(s, r.form)

		} else {

			exc.Write(s)

		}

		content := ""

		//获取一个随机结果(map的顺序不是固定的)，用做显示
		for _, s3 := range s {

			content = s3

			break
		}

		if r.form.Total != 0 {

			fmt.Print("当前进度：", decimal.NewFromInt(int64(r.form.CurrentIndex)).Div(decimal.NewFromInt(int64(r.form.Total))).Mul(decimal.NewFromInt(100)).String(), "%,", tools.SubStr(content, 0, conf.Conf.MaxStrLength)+"", "\r")

		}

	}

	if r.form.ResultCallback == nil {

		filename := exc.Save()

		r.form.Notice.PushMessage(notice.NewLog("excel文件为:" + filename))

	}

}
