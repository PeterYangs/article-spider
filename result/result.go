package result

import "github.com/PeterYangs/article-spider/form"

// GetResult 获取爬取结果
func GetResult(form form.Form, callback func(item map[string]string)) {

	defer form.StorageWait.Done()

	go checkChan(form)

	for v := range form.Storage {

		callback(v)

	}

}

//检查是否已完成爬取
func checkChan(form form.Form) {

	select {

	case <-form.IsFinish:

		//关闭通道写入
		close(form.Storage)

	}

}
