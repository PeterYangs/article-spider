package form

import (
	"article-spider/fileTypes"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/PeterYangs/tools"
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
	ListFields         map[string]Field       //列表页面字段选择器
	ExcelFile          *excelize.File         //excel表格对象
	Storage            chan map[string]string //存储爬取数据 ["title"]="文章标题"
	StorageTemp        map[string]string      //存储列表页数据
	ExcelWait          *sync.WaitGroup
	DetailMaxCoroutine int                    //爬取详情页最大协程数，默认按照列表的长度
	DisableAutoCoding  bool                   //是否关闭自动转码
	IsFinish           chan bool              //通知excel已完成爬取
	ProxyAddress       string                 //代理地址
	HttpHeader         map[string]string      //header
	HttpSetting        tools.HttpSetting      //全局http设置
	Uid                string                 //可视化下的websocket的uid
	BroadcastChan      chan map[string]string //广播管道
	CustomExcelHeader  bool                   //自定义Excel表格头部
	BroadcastWait      *sync.WaitGroup        //通知通道处理完毕等待
	DisableDebug       bool                   //是否关闭调试模式，开启调试模式后，所有的输出会在终端上
}

type Field struct {
	Types                fileTypes.FieldTypes
	Selector             string                   //选择器
	ImagePrefix          string                   //图片路径前缀,会生成到Excel表格中，但不会生成文件夹
	ImageDir             string                   //图片子文件夹，支持变量 1.[date:Y-m-d] 2.[random:1-100] 3.[singleField:title]
	ExcelHeader          string                   //excel表头，需要CustomExcelHeader为true,例：A
	ConversionFormatFunc func(data string) string //转换格式函数
	AttrKey              string                   //属性值参数
}
