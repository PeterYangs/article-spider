package result

import (
	"github.com/PeterYangs/article-spider/v2/excel"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/mode"
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

		//fmt.Println("准备退出")
		//
		//time.Sleep(10 * time.Second)
		//
		//fmt.Println("等待完毕")

		r.form.Wait.Done()

		r.form.Notice.Close()

	}()

	exc := excel.NewExcel(r.form)

	for s := range r.form.Storage {

		//合并列表和详情选择器
		var all = make(map[string]form.Field)

		tempRes := s

		for i, v := range r.form.ListFields {

			all[i] = v

		}

		for i, v := range r.form.DetailFields {

			all[i] = v

		}

		for i, v := range all {

			//自定义转换
			if v.ConversionFunc != nil {

				tempRes[i] = v.ConversionFunc(s[i], s)

			}

		}

		s = tempRes

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

		//fmt.Print("当前进度：", decimal.NewFromInt(int64(r.form.CurrentIndex)).Div(decimal.NewFromInt(int64(r.form.Total))).Mul(decimal.NewFromInt(100)).String(), "%,", tools.SubStr(content, 0, conf.Conf.MaxStrLength)+"", "\r")

		//fmt.Println("当前进度：", decimal.NewFromInt(int64(r.form.CurrentIndex)).Div(decimal.NewFromInt(int64(r.form.Total))).Mul(decimal.NewFromInt(100)).String(), "%,", tools.SubStr(content, 0, conf.Conf.MaxStrLength))

		switch r.form.Mode {

		case mode.Auto:

			r.form.Notice.Process("当前页码：", r.form.AutoPage+1, "/", r.form.Length, tools.SubStr(content, 0, r.form.Conf.MaxStrLength))
		default:

			if r.form.Total != 0 {

				r.form.Notice.Process("当前进度：", decimal.NewFromInt(int64(r.form.CurrentIndex)).Div(decimal.NewFromInt(int64(r.form.Total))).Mul(decimal.NewFromInt(100)).String(), "%,", tools.SubStr(content, 0, r.form.Conf.MaxStrLength))

			}

		}

	}

	if r.form.ResultCallback == nil {

		filename := exc.Save()

		//r.form.Notice.PushMessage(notice.NewLog("excel文件为:" + filename))

		r.form.Notice.Log("excel文件为:" + filename)

	}

}
