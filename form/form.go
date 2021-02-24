package form

import (
	"article-spider/fileTypes"
)

type Form struct {
	Host             string
	Channel          string
	Limit            int
	PageStart        int
	ListSelector     string
	ListHrefSelector string
	DetailFields     map[string]Field //详情页面字段选择器
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
