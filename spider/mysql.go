package spider

import (
	"article-spider/form"
	"fmt"
)

func WriteMysql(form form.Form) {

	//fmt.Println("gg111")

	defer form.MysqlWait.Done()

	go checkMysqlChan(form)

	for v := range form.Storage {

		//fmt.Println("gg")

		fmt.Println(v)

	}

}

func connect(form form.Form) {

}

func checkMysqlChan(form form.Form) {

	select {

	case <-form.IsFinish:

		//关闭通道写入
		close(form.Storage)

	}

}
