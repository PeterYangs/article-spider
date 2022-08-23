package article_spider

import (
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/shopspring/decimal"
	"strings"
	"sync"
)

type Rows struct {
	err  error
	maps map[string]string
}

func NewRows(m map[string]string) *Rows {

	//去空格
	for s, s2 := range m {

		m[s] = strings.TrimSpace(s2)
	}

	return &Rows{
		maps: m,
	}
}

// Add 合并两个结果
func (r *Rows) Add(rr *Rows) {

	if rr.err != nil {

		r.err = rr.err
	}

	for s, s2 := range rr.maps {

		r.maps[s] = strings.TrimSpace(s2)
	}

}

type result struct {
	s       *Spider
	storage chan *Rows
	lock    sync.Mutex
}

func NewResult(s *Spider) *result {

	return &result{s: s, storage: make(chan *Rows, 10), lock: sync.Mutex{}}
}

func (r *result) Push(i *Rows) {

	r.lock.Lock()

	defer r.lock.Unlock()

	select {

	case <-r.s.cxt.Done():

		return

	default:

		r.storage <- i
	}

}

func (r *result) Work() {

	r.s.wait.Add(1)

	defer func() {

		if r.s.form.ResultCallback == nil {

			filename := r.s.excel.Save()

			r.s.notice.Finish("excel文件为:" + "static/" + filename)

		} else {

			r.s.notice.Finish("输出完毕")

		}

		r.s.wait.Done()

	}()

	//优先结果管道
	for {

		select {

		case m := <-r.storage:

			if m.err != nil && r.s.form.FilterError {

				continue
			}

			r.do(m.maps)

		case <-r.s.cxt.Done():

			select {
			case m := <-r.storage:

				if m.err != nil && r.s.form.FilterError {

				}

				r.do(m.maps)

			default:

				fmt.Println()
				fmt.Println("结果协程退出")
				fmt.Println()

				return

			}

		}
	}

}

func (r *result) do(m map[string]string) {

	//合并列表和详情选择器
	var all = make(map[string]Field)

	tempRes := m

	for i, v := range r.s.form.ListFields {

		all[i] = v

	}

	for i, v := range r.s.form.DetailFields {

		all[i] = v

	}

	for i, v := range all {

		//自定义转换
		if v.ConversionFunc != nil {

			tempRes[i] = v.ConversionFunc(m[i], m)

		}

	}

	m = tempRes

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

		r.s.notice.Process("当前页码：", r.s.autoPage+1, "/", r.s.form.Length, strings.Replace(tools.SubStr(content, 0, 30), "\n", "", -1))
	default:

		if r.s.total != 0 {

			r.s.notice.Process("当前进度：", decimal.NewFromInt(int64(r.s.currentIndex)).Div(decimal.NewFromInt(int64(r.s.total))).Mul(decimal.NewFromInt(100)).String(), "%,", strings.Replace(tools.SubStr(content, 0, 30), "\n", "", -1))

		} else {

			r.s.notice.Process("正在计算进度")

		}

	}
}
