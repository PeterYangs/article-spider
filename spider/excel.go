package spider

import (
	"article-spider/form"
	"strconv"
)

//写入excel表
func WriteExcel(form form.Form) {

	defer form.ExcelWait.Done()

	row := 2

	array := make(map[string]string)

	index := 0

	var headerList []string

	//合并表头
	for i, _ := range form.DetailFields {

		headerList = append(headerList, i)
	}

	for i, _ := range form.ListFields {

		headerList = append(headerList, i)
	}

	//设置表头
	for _, i := range headerList {

		headName, arrayTemp := setHeader(i, index, array)

		array = arrayTemp

		form.ExcelFile.SetCellValue("Sheet1", headName+"1", i)

		index++

	}

	go checkChan(form)

	//写入
	for v := range form.Storage {

		indexs := 0

		for ii, vv := range v {

			if indexs == 0 {

				//fmt.Println(vv)

				//输出到通知管道
				form.BroadcastChan <- map[string]string{"types": "log", "data": vv}

			}

			form.ExcelFile.SetCellValue("Sheet1", getHeader(ii, array)+strconv.Itoa(row), vv)

			indexs++
		}

		row++

	}

}

// array ["title"]="A"
func setHeader(name string, index int, array map[string]string) (string, map[string]string) {

	headerList := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA"}

	head := headerList[index]

	array[name] = head

	return head, array
}

func getHeader(name string, array map[string]string) string {

	return array[name]
}

//检查是否已完成爬取
func checkChan(form form.Form) {

	select {

	case <-form.IsFinish:

		//关闭通道写入
		close(form.Storage)

	}

}
