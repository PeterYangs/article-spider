package spider

import "github.com/PeterYangs/article-spider/v2/form"

type Spider struct {
	form *form.Form
}

func NewSpider() *Spider {

	return &Spider{}
}

func (s *Spider) LoadForm(form *form.Form) *Spider {

	s.form = form

	return s
}

func (s *Spider) Start() {

}
