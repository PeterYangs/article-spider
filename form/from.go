package form

import (
	"errors"
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/mode"
	"github.com/PeterYangs/article-spider/v2/notice"
	"github.com/PeterYangs/request"
	"github.com/PeterYangs/tools"

	"github.com/PuerkitoBio/goquery"
	uuid "github.com/satori/go.uuid"
	http2 "net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CustomForm struct {
	Host                       string                              //网站域名
	Channel                    string                              //栏目链接，页码用[PAGE]替换
	PageStart                  int                                 //页码起始页
	Length                     int                                 //爬取页码长度
	ListSelector               string                              //列表选择器
	HrefSelector               string                              //a链接选择器，相对于列表选择器
	DisableAutoCoding          bool                                //是否自动转码
	LazyImageAttrName          string                              //懒加载图片属性，默认为data-original
	DisableImageExtensionCheck bool                                //禁用图片拓展名检查，禁用后所有图片拓展名强制为png
	AllowImageExtension        []string                            //允许下载的图片拓展名
	DefaultImg                 func(form *Form, item Field) string //图片出错时，设置默认图片
	DetailFields               map[string]Field                    //详情页面字段选择器
	ListFields                 map[string]Field                    //列表页面字段选择器,暂不支持api爬取
	CustomExcelHeader          bool                                //自定义Excel表格头部
	DetailCoroutineNumber      int                                 //爬取详情页协程数
	HttpTimeout                time.Duration                       //请求超时时间
	HttpHeader                 map[string]string                   //header
}

type Form struct {
	Host                       string          //网站域名
	Channel                    string          //栏目链接，页码用[PAGE]替换
	PageStart                  int             //页码起始页
	Length                     int             //爬取页码长度
	Client                     *request.Client //http客户端
	ListSelector               string          //列表选择器
	HrefSelector               string          //a链接选择器，相对于列表选择器
	Mode                       mode.Mode
	DisableAutoCoding          bool //是否自动转码
	Notice                     *notice.Notice
	Wait                       sync.WaitGroup
	LazyImageAttrName          string                              //懒加载图片属性，默认为data-original
	DisableImageExtensionCheck bool                                //禁用图片拓展名检查，禁用后所有图片拓展名强制为png
	AllowImageExtension        []string                            //允许下载的图片拓展名
	DefaultImg                 func(form *Form, item Field) string //图片出错时，设置默认图片
	DetailFields               map[string]Field                    //详情页面字段选择器
	ListFields                 map[string]Field                    //列表页面字段选择器,暂不支持api爬取
	Storage                    chan map[string]string              //数据结果通道
	CustomExcelHeader          bool                                //自定义Excel表格头部
	DetailCoroutineNumber      int                                 //爬取详情页协程数
	DetailCoroutineChan        chan bool                           //限制详情页并发chan
	DetailWait                 sync.WaitGroup
	HttpTimeout                time.Duration     //请求超时时间
	HttpHeader                 map[string]string //header
	DetailSize                 int               //每个列表的详情数量
	Total                      int               //预计爬取总数
	CurrentIndex               int               //当前爬取数量
}

type Field struct {
	Types       fileTypes.FieldTypes
	Selector    string                               //字段选择器
	AttrKey     string                               //属性值参数
	ImagePrefix func(form *Form, path string) string //图片路径前缀,会添加到图片路径前缀，但不会生成文件夹
	ImageDir    string                               //图片子文件夹，支持变量 1.[date:Y-m-d] 2.[random:1-100] 3.[singleField:title]
	ExcelHeader string                               //excel表头，需要CustomExcelHeader为true,例：A
}

// DealCoding 解决编码问题
func (f *Form) DealCoding(html string, header http2.Header) (string, error) {

	headerContentType_ := header["Content-Type"]

	if len(headerContentType_) > 0 {

		headerContentType := headerContentType_[0]

		charset := f.GetCharsetByContentType(headerContentType)

		switch charset {

		case "gbk":

			return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

		case "gb2312":

			return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

		case "utf-8":

			return html, nil

		case "utf8":

			return html, nil

		}

	}

	code, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {

		return html, err
	}

	contentType, _ := code.Find("meta[charset]").Attr("charset")

	//转小写
	contentType = strings.TrimSpace(strings.ToLower(contentType))

	switch contentType {

	case "gbk":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	case "gb2312":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	case "utf-8":

		return html, nil

	case "utf8":

		return html, nil

	}

	contentType, _ = code.Find("meta[http-equiv=\"Content-Type\"]").Attr("content")

	charset := f.GetCharsetByContentType(contentType)

	switch charset {

	case "gbk":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	case "gb2312":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	}

	return html, nil
}

// GetCharsetByContentType 从contentType中获取编码
func (f *Form) GetCharsetByContentType(contentType string) string {

	contentType = strings.TrimSpace(strings.ToLower(contentType))

	//捕获编码
	r, _ := regexp.Compile(`charset=([^;]+)`)

	re := r.FindAllStringSubmatch(contentType, 1)

	if len(re) > 0 {

		c := re[0][1]

		return c

	}

	return ""
}

// GetHref 获取完整a链接
func (f *Form) GetHref(href string) string {

	case1, _ := regexp.MatchString("^/[a-zA-Z0-9_]+.*", href)

	case2, _ := regexp.MatchString("^//[a-zA-Z0-9_]+.*", href)

	case3, _ := regexp.MatchString("^(http|https).*", href)

	switch true {

	case case1:

		href = f.Host + href

		break

	case case2:

		//获取当前网址的协议
		res := regexp.MustCompile("^(https|http).*").FindStringSubmatch(f.Host)

		href = res[1] + ":" + href

		break

	case case3:

		break

	default:

		href = f.Host + "/" + href
	}

	return href

}

// ResolveSelector 解析选择器
func (f *Form) ResolveSelector(html string, selector map[string]Field, originUrl string) (map[string]string, error) {

	//存储结果
	var res = &sync.Map{}

	var wait = &sync.WaitGroup{}

	//goquery加载html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {

		return nil, err

	}

	//解析详情页面选择器
	for field, item := range selector {

		switch item.Types {

		//单个文字字段
		case fileTypes.Text:

			v := doc.Find(item.Selector).Text()

			res.Store(field, v)

			break

		//单个文字字段
		case fileTypes.Attr:

			v, _ := doc.Find(item.Selector).Attr(item.AttrKey)

			res.Store(field, v)

			break

		//只爬html（不包括图片）
		case fileTypes.OnlyHtml:

			v, err := doc.Find(item.Selector).Html()

			if err != nil {

				res.Store(field, "")

				f.Notice.PushMessage(notice.NewError(err.Error()+",源链接："+originUrl, ",选择器：", item.Selector))

				break

			}

			res.Store(field, v)

			break

		//爬取html，包括图片
		case fileTypes.HtmlWithImage:

			wait.Add(1)

			go func(item Field, field string) {

				defer wait.Done()

				html_, err := doc.Find(item.Selector).Html()

				if err != nil {

					f.Notice.PushMessage(notice.NewError(err.Error()+",源链接："+originUrl, ",选择器：", item.Selector))

					return

				}

				htmlImg, err := goquery.NewDocumentFromReader(strings.NewReader(html_))

				if err != nil {

					f.Notice.PushMessage(notice.NewError(err.Error() + ",源链接：" + originUrl))

					return

				}

				var waitImg sync.WaitGroup

				var imgList = sync.Map{}

				htmlImg.Find("img").Each(func(i int, selection *goquery.Selection) {

					img, err := f.getImageLink(selection)

					if err != nil {

						f.Notice.PushMessage(notice.NewError(err.Error()+",源链接："+originUrl, ",富文本内容"))

						return
					}

					waitImg.Add(1)

					go func(waitImg *sync.WaitGroup, imgList *sync.Map) {

						defer waitImg.Done()

						imgName := f.DownImg(img, item, res)

						imgList.Store(imgName, img)

					}(&waitImg, &imgList)

				})

				waitImg.Wait()

				imgList.Range(func(key, value interface{}) bool {

					html_ = strings.Replace(html_, value.(string), key.(string), -1)

					return true
				})

				res.Store(field, html_)

			}(item, field)

		//单个图片
		case fileTypes.Image:

			wait.Add(1)

			go func(item Field, field string) {

				defer wait.Done()

				imgUrl, err := f.getImageLink(doc.Find(item.Selector))

				if err != nil {

					f.Notice.PushMessage(notice.NewError(err.Error()+",源链接："+originUrl, ",选择器：", item.Selector))

					return
				}

				imgName := f.DownImg(imgUrl, item, res)

				res.Store(field, imgName)

			}(item, field)

			break

		//多个图片
		case fileTypes.MultipleImages:

			wait.Add(1)

			go func(_item Field, field string) {

				defer wait.Done()

				var waitImg sync.WaitGroup

				var imgList = sync.Map{}

				//sync.

				doc.Find(item.Selector).Each(func(i int, selection *goquery.Selection) {

					imgUrl, err := f.getImageLink(selection)

					if err != nil {

						f.Notice.PushMessage(notice.NewError(err.Error()+",源链接："+originUrl, ",选择器：", item.Selector))

						return
					}

					waitImg.Add(1)

					go func(waitImg *sync.WaitGroup, imgList *sync.Map) {

						defer waitImg.Done()

						imgName := f.DownImg(imgUrl, item, res)

						imgList.Store(imgName, "")

					}(&waitImg, &imgList)

				})

				waitImg.Wait()

				var strArray []string

				imgList.Range(func(key, value interface{}) bool {

					strArray = append(strArray, key.(string))

					return true
				})

				array := tools.Join(",", strArray)

				res.Store(field, array)

			}(item, field)

		//固定数据
		case fileTypes.Fixed:

			res.Store(field, item.Selector)

		}

	}

	wait.Wait()

	arr := make(map[string]string)

	res.Range(func(key, value interface{}) bool {

		arr[key.(string)] = value.(string)

		return true

	})

	return arr, nil

}

//获取图片链接
func (f *Form) getImageLink(imageDoc *goquery.Selection) (string, error) {

	if f.LazyImageAttrName != "" {

		imgUrl, imgBool := imageDoc.Attr(f.LazyImageAttrName)

		if imgBool == false || imgUrl == "" {

			imgUrl, imgBool = imageDoc.Attr("src")

			if imgBool == false || imgUrl == "" {

				return "", errors.New("未找到图片链接")
			}

		}

		return imgUrl, nil
	}

	imgUrl, imgBool := imageDoc.Attr("src")

	if imgBool == false || imgUrl == "" {

		//懒加载
		imgUrl, imgBool = imageDoc.Attr("data-original")

		if imgBool == false || imgUrl == "" {

			return "", errors.New("未找到图片链接")
		}

	}

	return imgUrl, nil
}

// DownImg 下载图片（包括生成文件夹）
func (f *Form) DownImg(url string, item Field, res *sync.Map) string {

	//获取完整链接
	imgUrl := f.GetHref(url)

	//生成随机名称
	uuidString := uuid.NewV4().String()

	uuidString = strings.Replace(uuidString, "-", "", -1)

	dir := ""

	if item.ImageDir != "" {

		//获取图片文件夹
		dir = f.GetDir(item.ImageDir, res)

		//panic(dir)

		//设置文件夹
		err := tools.MkDirDepth("image/" + dir)

		if err != nil {

			//ErrorLine(form, err.Error())

			f.Notice.PushMessage(notice.NewError(err.Error()))

			//return
		}
	}

	ex, err := tools.GetExtensionName(imgUrl)

	if err != nil {

		ex = "png"
	}

	//禁用拓展名检查
	if f.DisableImageExtensionCheck {

		ex = "png"

	} else {

		allowImage := []string{"png", "jpg", "jpeg", "gif", "jfif"}

		//自定义允许下载的图片拓展名
		if len(f.AllowImageExtension) > 0 {

			allowImage = f.AllowImageExtension
		}

		if !tools.In_array(allowImage, strings.ToLower(ex)) {

			//ErrorLine(form, "图片拓展名异常:"+imgUrl)

			f.Notice.PushMessage(notice.NewError("图片拓展名异常:" + imgUrl))

			//获取默认图片
			if f.DefaultImg != nil {

				return f.DefaultImg(f, item)
			}

			return ""
		}

	}

	imgName := (If(dir == "", "", dir+"/")).(string) + uuidString + "." + ex

	//imgErr := f.Client.Request().DownloadFile(imgUrl, "image/"+imgName)

	imgErr := f.Client.R().Download(imgUrl, "image/"+imgName)

	if imgErr != nil {

		msg := imgErr.Error()

		//ErrorLine(form, msg)

		f.Notice.PushMessage(notice.NewError(msg))

		//获取默认图片
		if f.DefaultImg != nil {

			return f.DefaultImg(f, item)
		}

		return ""

	}

	prefix := ""

	if item.ImagePrefix != nil {

		prefix = item.ImagePrefix(f, imgName)

	}

	//自动添加斜杠
	if tools.SubStr(prefix, -1, -1) != "/" {

		prefix += "/"
	}

	return (If(item.ImagePrefix == nil, "", prefix)).(string) + imgName

}

func (f *Form) GetDir(path string, res *sync.Map) string {

	//替换时间格式
	r1, _ := regexp.Compile(`\[date:(.*?)]`)

	date := r1.FindAllStringSubmatch(path, -1)

	for _, v := range date {

		path = strings.Replace(path, v[0], tools.Date(v[1], time.Now().Unix()), -1)

	}

	//替换随机格式
	r2, _ := regexp.Compile(`\[random:([0-9]+-[0-9]+)]`)

	random := r2.FindAllStringSubmatch(path, -1)

	for _, v := range random {

		min, _ := strconv.Atoi(tools.Explode("-", v[1])[0])

		max, _ := strconv.Atoi(tools.Explode("-", v[1])[1])

		path = strings.Replace(path, v[0], strconv.FormatInt(tools.Mt_rand(int64(min), int64(max)), 10), -1)

	}

	//根据爬取文件给文件夹命名
	r3, _ := regexp.Compile(`\[singleField:(.*?)]`)

	singleField := r3.FindAllStringSubmatch(path, -1)

	for i, v := range singleField {

		field := ""

		//ok:=false

		if i == 0 {

			times := 0

			for {

				field_, ok := res.Load(v[1])

				if !ok {

					time.Sleep(200 * time.Millisecond)

					times++

					if times >= 5 {

						field = "timeout"

						break
					}

				} else {

					field = field_.(string)

					//处理为空的情况
					if field == "" {

						field = "unknown"
					}

					break

				}

			}

		}

		path = strings.Replace(path, v[0], field, -1)

	}

	return path

}

// If 伪三元运算
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
