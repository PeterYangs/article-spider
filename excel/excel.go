package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PeterYangs/article-spider/v3/form"
	ff "github.com/PeterYangs/article-spider/v3/form"
	uuid "github.com/satori/go.uuid"
	"os"
	"strconv"
)

type excel struct {
	file        *excelize.File
	form        *form.Form
	headerArray map[string]string
	line        int
}

func NewExcel(form *form.Form) *excel {

	f := excelize.NewFile()

	f.SetActiveSheet(f.NewSheet("Sheet1"))

	e := &excel{file: f, form: form}

	e.line = 2

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

		headName, arrayTemp := e.setHeader(i, index, array, v)

		array = arrayTemp

		e.file.SetCellValue("Sheet1", headName+"1", i)

		index++

	}

	e.headerArray = array

	return e
}

func (e *excel) Write(item map[string]string) {

	for ii, vv := range item {

		e.file.SetCellValue("Sheet1", e.getHeader(ii, e.headerArray)+strconv.Itoa(e.line), vv)

	}

	e.line++

}

func (e *excel) Save() string {

	os.Mkdir("static", 0755)

	filename := uuid.NewV4().String()

	e.file.SaveAs("static/" + filename + ".xlsx")

	return filename

}

// array ["title"]="A"
func (e *excel) setHeader(name string, index int, array map[string]string, item ff.Field) (string, map[string]string) {

	headerList := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA"}

	if e.form.CustomExcelHeader {

		head := item.ExcelHeader

		array[name] = head

		return head, array

	} else {

		head := headerList[index]

		array[name] = head

		return head, array

	}

}

func (e *excel) getHeader(name string, array map[string]string) string {

	return array[name]
}
