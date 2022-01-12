package article_spider

import (
	"github.com/PeterYangs/tools"
	"github.com/shopspring/decimal"
)

type result struct {
	s       *Spider
	storage chan map[string]string
}

func NewResult(s *Spider) *result {

	return &result{s: s, storage: make(chan map[string]string, 10)}
}

func (r *result) Push(m map[string]string) {

	r.storage <- m

}

func (r *result) Work() {

	r.s.wait.Add(1)

	defer func() {

		if r.s.form.ResultCallback == nil {

			filename := r.s.excel.Save()

			//r.form.Notice.Log("excel文件为:" + filename)

			r.s.notice.Finish("excel文件为:" + filename)

		}

		r.s.wait.Done()

	}()

	for {

		select {

		case m := <-r.storage:

			//自定义结果处理
			if r.s.form.ResultCallback != nil {

				r.s.form.ResultCallback(m, &r.s.form)

			} else {

				r.s.excel.Write(m)

			}

			content := ""

			//获取一个随机结果(map的顺序不是固定的)，用做显示
			for _, s3 := range m {

				content = s3

				break
			}

			//fmt.Println(content)

			switch r.s.mode {

			case Auto:

				r.s.notice.Process("当前页码：", r.s.autoPage+1, "/", r.s.form.Length, tools.SubStr(content, 0, 30))
			default:

				if r.s.total != 0 {

					r.s.notice.Process("当前进度：", decimal.NewFromInt(int64(r.s.currentIndex)).Div(decimal.NewFromInt(int64(r.s.total))).Mul(decimal.NewFromInt(100)).String(), "%,", tools.SubStr(content, 0, 30))

				} else {

					r.s.notice.Process("正在计算进度")

				}

			}

		case <-r.s.cxt.Done():

			return

		}
	}

}
