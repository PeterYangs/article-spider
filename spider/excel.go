package spider

import (
	"github.com/PeterYangs/article-spider/form"
	ff "github.com/PeterYangs/article-spider/form"
	"strconv"
)

// WriteExcel 写入excel表
func WriteExcel(form form.Form) {

	defer form.ExcelWait.Done()

	row := 2

	array := make(map[string]string)

	index := 0

	var headerList = make(map[string]ff.Field)

	//合并表头
	for i, v := range form.DetailFields {

		headerList[i] = v
	}

	for i, v := range form.ListFields {

		headerList[i] = v
	}

	//设置表头
	for i, v := range headerList {

		headName, arrayTemp := setHeader(i, index, array, v, form)

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
func setHeader(name string, index int, array map[string]string, item ff.Field, form ff.Form) (string, map[string]string) {

	headerList := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA"}

	if form.CustomExcelHeader {

		head := item.ExcelHeader

		array[name] = head

		return head, array

	} else {

		head := headerList[index]

		array[name] = head

		return head, array

	}

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
