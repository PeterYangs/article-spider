package form

import (
	"article-spider/fileTypes"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"sync"
)

type Form struct {
	Host               string
	Channel            string
	Limit              int
	PageStart          int
	ListSelector       string
	ListHrefSelector   string
	DetailFields       map[string]Field       //详情页面字段选择器
	ExcelFile          *excelize.File         //excel表格对象
	Storage            chan map[string]string //存储爬取数据 ["title"]="文章标题"
	ExcelWait          *sync.WaitGroup
	DetailMaxCoroutine int //爬取详情页最大协程数，默认按照列表的长度

}

type Field struct {
	Types          fileTypes.FieldTypes
	SingleSelector string
}

func getForm() {

	//y:=FieldTypes.

	//FieldTypes.

	//os.OpenFile("", os.O_CREATE, 0777)

	//f:=Field{
	//	fileTypes.SingleField,
	//}

}
