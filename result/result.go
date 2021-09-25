package result

import (
	"github.com/PeterYangs/article-spider/v2/excel"
	"github.com/PeterYangs/article-spider/v2/form"
	"github.com/PeterYangs/article-spider/v2/notice"
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

		exc.Write(s)

		r.form.Notice.PushMessage(notice.NewInfo(s))

	}

	filename := exc.Save()

	r.form.Notice.PushMessage(notice.NewInfo("excel文件为:" + filename))

}
