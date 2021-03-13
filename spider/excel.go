package spider

import (
	"article-spider/form"
	"strconv"
)

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

	//写入
	for v := range form.Storage {

		//fmt.Println(v)

		for ii, vv := range v {

			form.ExcelFile.SetCellValue("Sheet1", getHeader(ii, array)+strconv.Itoa(row), vv)
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
