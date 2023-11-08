package article_spider

type FieldTypes int

const (
	Text           FieldTypes = 0x00000 //单个字段
	Image          FieldTypes = 0x00002 //单个图片
	OnlyHtml       FieldTypes = 0x00003 //普通html(不包括图片)
	HtmlWithImage  FieldTypes = 0x00004 //html包括图片
	MultipleImages FieldTypes = 0x00005 //多图
	Attr           FieldTypes = 0x00006 //标签属性选择器
	Fixed          FieldTypes = 0x00007 //固定数据，填什么返回什么,选择器就是返回的数据
	Regular        FieldTypes = 0x00008 //正则（FindStringSubmatch,返回一个结果）
	File           FieldTypes = 0x00009 //文件类型
	Attrs          FieldTypes = 0x00010 //属性列表，如一个图片列表的所有图片链接
)
