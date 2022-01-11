package article_spider

import "github.com/PeterYangs/request"

type Spider struct {
	form   Form
	mode   Mode
	client *request.Client
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

	return &Spider{form: f, mode: mode, client: client}
}

func (s *Spider) Start() {

	switch s.mode {

	case Normal:

		NewNormal(s).Start()

	}

}
