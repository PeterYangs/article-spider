package article_spider

import "time"

type Form struct {
	Host              string            //网站域名
	Channel           string            //栏目链接，页码用[PAGE]替换
	PageStart         int               //页码起始页
	Length            int               //爬取页码长度
	ListSelector      string            //列表选择器
	HrefSelector      string            //a链接选择器，相对于列表选择器
	DisableAutoCoding bool              //是否禁用自动转码
	DetailFields      map[string]Field  //详情页面字段选择器
	ListFields        map[string]Field  //列表页面字段选择器,暂不支持api爬取
	HttpTimeout       time.Duration     //请求超时时间
	HttpHeader        map[string]string //header
	HttpProxy         string            //代理（暂不支持auto模式，但是下载图片只有的）

}

type Field struct {
	Types        FieldTypes
	Selector     string                                    //字段选择器
	AttrKey      string                                    //属性值参数
	ImagePrefix  func(form *Form, imageName string) string //图片路径前缀,会添加到图片路径前缀，但不会生成文件夹
	ImageDir     string                                    //图片子文件夹，支持变量 1.[date:Y-m-d] 2.[random:1-100] 3.[singleField:title]
	ExcelHeader  string                                    //excel表头，需要CustomExcelHeader为true,例：A
	RegularIndex int                                       //正则匹配中的反向引用的下标，默认是1
}
