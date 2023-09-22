package article_spider

import (
	"context"
	"errors"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"github.com/phpisfirstofworld/image"
	uuid "github.com/satori/go.uuid"
	http2 "net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Form struct {
	Host                       string                                   //网站域名
	Channel                    string                                   //栏目链接，页码用[PAGE]替换
	PageStart                  int                                      //页码起始页
	PageCurrent                int                                      //当前页码
	ListUrlCurrent             string                                   //当前列表链接
	Length                     int                                      //爬取页码长度
	ListSelector               string                                   //列表选择器
	HrefSelector               string                                   //a链接选择器，相对于列表选择器
	DisableAutoCoding          bool                                     //是否禁用自动转码
	DetailFields               map[string]Field                         //详情页面字段选择器
	ListFields                 map[string]Field                         //列表页面字段选择器,暂不支持api爬取
	HttpTimeout                time.Duration                            //请求超时时间
	HttpHeader                 map[string]string                        //header
	HttpProxy                  string                                   //代理
	ChannelFunc                func(form *Form) []string                //自定义栏目链接
	DetailCoroutineNumber      int                                      //爬取详情页协程数
	LazyImageAttrName          string                                   //懒加载图片属性，默认为data-original
	DisableImageExtensionCheck bool                                     //禁用图片拓展名检查，禁用后所有图片拓展名强制为png
	AllowImageExtension        []string                                 //允许下载的图片拓展名
	DefaultImg                 func(form *Form, item Field) string      //图片出错时，设置默认图片
	MiddleSelector             []string                                 //中间层选择器(a链接选择器)，当详情页有多层时使用，暂不支持自动模式
	CustomExcelHeader          bool                                     //自定义Excel表格头部
	ResultCallback             func(item map[string]string, form *Form) //自定义获取爬取结果回调
	ApiConversion              func(html string, form *Form) []string   //api获取链接
	AutoPrefixEvent            func(chromedpCtx context.Context)        //自动爬取模式前置事件
	AutoListWaitSelector       string                                   //列表等待选择器（用于自动化爬取）
	AutoNextPageMode           NextPageMode                             //下一页模式（用于自动化爬取,目前支持常规分页和加载更多）
	AutoDetailForceNewTab      bool                                     //自动模式详情页强制打开新窗口(必须是a链接)
	AutoDetailWaitSelector     string                                   //详情等待选择器（用于自动化爬取）
	AutoNextSelector           string                                   //下一页选择器（用于自动化爬取）
	FilterError                bool                                     //过滤错误的行
	DetailUrls                 []string                                 //详情页列表
	Filter                     func(map[string]string) bool             //数据过滤，返回false则放弃数据
	s                          *Spider
}

type Field struct {
	Types              FieldTypes
	Selector           string                                              //字段选择器
	NotSelector        []string                                            //剔除选择器(后置选择器，意思是先获取该item的doc再剔除节点)
	PrefixNotSelector  []string                                            //前置剔除选择器(意思是先剔除html的节点)
	AttrKey            string                                              //属性值参数
	ImagePrefix        func(form *Form, imageName string) string           //图片路径前缀,会添加到图片路径前缀，但不会生成文件夹
	ImageDir           string                                              //图片子文件夹，支持变量 1.[date:Y-m-d] 2.[random:1-100] 3.[singleField:title]
	ExcelHeader        string                                              //excel表头，需要CustomExcelHeader为true,例：A
	RegularIndex       int                                                 //正则匹配中的反向引用的下标，默认是1
	ConversionFunc     func(data string, resList map[string]string) string //转换格式函数,第一个参数是该字段数据，第二个参数是所有数据，跟web框架的获取器类似
	LazyImageAttrName  string                                              //懒加载图片属性，默认为data-original
	ImageResizePercent int                                                 //图片缩放百分比

}

// DealCoding 解决编码问题
func (f *Form) DealCoding(html string, header http2.Header) (string, error) {

	headerContentType_ := header["Content-Type"]

	if len(headerContentType_) > 0 {

		headerContentType := headerContentType_[0]

		charset := f.GetCharsetByContentType(headerContentType)

		charset = strings.ToLower(charset)

		switch charset {

		case "gbk":

			return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

		case "gb2312":

			return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

		case "utf-8":

			return html, nil

		case "utf8":

			return html, nil

		case "euc-jp":

			return string(tools.ConvertToByte(html, "euc-jp", "utf8")), nil

		case "":

			break

		default:
			return string(tools.ConvertToByte(html, charset, "utf8")), nil

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

	case "euc-jp":

		return string(tools.ConvertToByte(html, "euc-jp", "utf8")), nil

	case "":

		break
	default:
		return string(tools.ConvertToByte(html, contentType, "utf8")), nil

	}

	contentType, _ = code.Find("meta[http-equiv=\"Content-Type\"]").Attr("content")

	charset := f.GetCharsetByContentType(contentType)

	switch charset {

	case "utf-8":

		return html, nil

	case "utf8":

		return html, nil

	case "gbk":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	case "gb2312":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	case "euc-jp":

		return string(tools.ConvertToByte(html, "euc-jp", "utf8")), nil

	case "":

		break

	default:
		return string(tools.ConvertToByte(html, charset, "utf8")), nil

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
func (f *Form) ResolveSelector(html string, selector map[string]Field, originUrl string) (*Rows, error) {

	//存储结果
	var res = &sync.Map{}

	var wait = &sync.WaitGroup{}

	var globalErr error = nil

	//goquery加载html
	htmlDoc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {

		return nil, err

	}

	//解析详情页面选择器
	for fieldT, itemT := range selector {

		doc := htmlDoc

		//前置剔除选择器
		for _, s := range itemT.PrefixNotSelector {

			doc.Find(s).Remove()
		}

		field := fieldT

		item := itemT

		switch item.Types {

		//单个文字字段
		case Text:

			selectors := doc.Find(item.Selector)

			//排除选择器
			for _, s := range item.NotSelector {

				selectors.Find(s).Remove()

			}

			v := strings.TrimSpace(selectors.Text())

			res.Store(field, v)

			break

		//单个元素属性
		case Attr:

			v := ""

			if strings.TrimSpace(item.Selector) == "" {

				v, _ = doc.Attr(item.AttrKey)

			} else {

				v, _ = doc.Find(item.Selector).Attr(item.AttrKey)
			}

			res.Store(field, strings.TrimSpace(v))

			break

		//只爬html（不包括图片）
		case OnlyHtml:

			//v, err := doc.Find(item.Selector).Html()

			selectors := doc.Find(item.Selector)

			//排除选择器
			for _, s := range item.NotSelector {

				selectors.Find(s).Remove()

			}

			v, sErr := selectors.Html()

			if sErr != nil {

				res.Store(field, "")

				f.s.notice.Error(err.Error()+",源链接："+originUrl, ",选择器：", item.Selector)

				globalErr = err

				break

			}

			res.Store(field, v)

			break

		//爬取html，包括图片
		case HtmlWithImage:

			wait.Add(1)

			go func(_item Field, field string) {

				defer wait.Done()

				selectors := doc.Find(_item.Selector)

				//排除选择器
				for _, s := range item.NotSelector {

					selectors.Find(s).Remove()

				}

				html_, sErr := selectors.Html()

				if sErr != nil {

					f.s.notice.Error(sErr.Error()+",源链接："+originUrl, ",选择器：", _item.Selector)

					globalErr = sErr

					return

				}

				htmlImg, err := goquery.NewDocumentFromReader(strings.NewReader(html_))

				if err != nil {

					f.s.notice.Error(err.Error() + ",源链接：" + originUrl)

					globalErr = err

					return

				}

				var waitImg sync.WaitGroup

				var imgList = sync.Map{}

				htmlImg.Find("img").Each(func(i int, selection *goquery.Selection) {

					img, err := f.getImageLink(selection, _item, originUrl)

					if err != nil {

						f.s.notice.Error(err.Error()+",源链接："+originUrl, ",富文本内容")

						globalErr = err

						return
					}

					waitImg.Add(1)

					go func(waitImg *sync.WaitGroup, imgList *sync.Map, __item Field) {

						defer waitImg.Done()

						imgName, e := f.DownImg(img, __item, res)

						if e != nil {

							f.s.notice.Error(e.Error()+",源链接："+originUrl, ",富文本图片下载失败", "图片地址", img)

						}

						globalErr = e

						imgList.Store(imgName, img)

					}(&waitImg, &imgList, _item)

				})

				waitImg.Wait()

				html_, _ = htmlImg.Html()

				imgList.Range(func(key, value interface{}) bool {

					html_ = strings.Replace(html_, value.(string), key.(string), -1)

					return true
				})

				res.Store(field, html_)

			}(item, field)

		//单个图片
		case Image:

			wait.Add(1)

			go func(_item Field, field string) {

				defer wait.Done()

				imgUrl, err := f.getImageLink(doc.Find(_item.Selector), _item, originUrl)

				if err != nil {

					f.s.notice.Error(err.Error()+",源链接："+originUrl, ",选择器：", _item.Selector)

					globalErr = err

					return
				}

				imgName, e := f.DownImg(imgUrl, _item, res)

				globalErr = e

				if e != nil {

					f.s.notice.Error(e.Error()+",源链接："+originUrl, ",选择器：", _item.Selector, "图片地址", imgUrl)
				}

				res.Store(field, imgName)

			}(item, field)

			break

		//单个文件
		case File:

			selectors := doc.Find(item.Selector)

			v, ok := selectors.Attr(item.AttrKey)

			if !ok {

				break
			}

			imgName, e := f.DownImg(v, item, res)

			globalErr = e

			res.Store(field, imgName)

			//res.Store(field, v)

		//多个图片
		case MultipleImages:

			wait.Add(1)

			go func(_item Field, field string) {

				defer wait.Done()

				var waitImg sync.WaitGroup

				var imgList = sync.Map{}

				doc.Find(_item.Selector).Each(func(i int, selection *goquery.Selection) {

					imgUrl, err := f.getImageLink(selection, _item, originUrl)

					if err != nil {

						f.s.notice.Error(err.Error()+",源链接："+originUrl, ",选择器：", _item.Selector)

						globalErr = err

						return
					}

					waitImg.Add(1)

					go func(waitImg *sync.WaitGroup, imgList *sync.Map, __item Field) {

						defer waitImg.Done()

						imgName, e := f.DownImg(imgUrl, __item, res)

						if e != nil {

							f.s.notice.Error(e.Error()+",源链接："+originUrl, ",选择器：", _item.Selector, "图片地址", imgUrl)

						}

						globalErr = e

						imgList.Store(imgName, "")

					}(&waitImg, &imgList, _item)

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
		case Fixed:

			res.Store(field, item.Selector)

		//正则
		case Regular:

			reg := regexp.MustCompile(item.Selector).FindStringSubmatch(html)

			if len(reg) > 0 {

				index := 1

				if item.RegularIndex != 0 {

					index = item.RegularIndex
				}

				res.Store(field, reg[index])

			}

			globalErr = errors.New("正则匹配未找到")

			f.s.notice.Error("正则匹配未找到")

		}

	}

	wait.Wait()

	arr := make(map[string]string)

	res.Range(func(key, value interface{}) bool {

		arr[key.(string)] = value.(string)

		return true

	})

	r := NewRows(arr)

	r.err = globalErr

	return r, nil

}

//获取图片链接
func (f *Form) getImageLink(imageDoc *goquery.Selection, item Field, originUrl string) (string, error) {

	//懒加载图片处理
	if item.LazyImageAttrName != "" {

		//Field里面的懒加载属性
		imgUrl, imgBool := imageDoc.Attr(item.LazyImageAttrName)

		if imgBool && imgUrl != "" {

			//填充图片src，防止图片无法显示
			imageDoc.RemoveAttr(item.LazyImageAttrName)

			imageDoc.SetAttr("src", imgUrl)

			return imgUrl, nil
		}

	}

	//懒加载图片处理
	if f.LazyImageAttrName != "" {

		//form里面的懒加载属性
		imgUrl, imgBool := imageDoc.Attr(f.LazyImageAttrName)

		if imgBool && imgUrl != "" {

			//填充图片src，防止图片无法显示
			imageDoc.RemoveAttr(f.LazyImageAttrName)

			imageDoc.SetAttr("src", imgUrl)

			return imgUrl, nil
		}

	}

	imgUrl, imgBool := imageDoc.Attr("src")

	if imgBool == false || imgUrl == "" {

		return "", errors.New("未找到图片链接，请检查是否存在懒加载")
	}

	return imgUrl, nil
}

func (f *Form) completePath(path string) string {

	if path == "" {

		return path
	}

	m, _ := regexp.MatchString(`.*/$`, path)

	if m {

		return path
	}

	return path + "/"
}

// DownImg 下载图片（包括生成文件夹）
func (f *Form) DownImg(url string, item Field, res *sync.Map) (string, error) {

	url = strings.Replace(url, "\n", "", -1)

	//获取完整链接
	imgUrl := f.GetHref(url)

	//生成随机名称
	uuidString := uuid.NewV4().String()

	uuidString = strings.Replace(uuidString, "-", "", -1)

	dir := ""

	//获取图片文件夹
	dir = f.GetDir(item.ImageDir, res)

	//设置文件夹,图片保存路径+图片默认前缀路径+生成路径
	err := os.MkdirAll(f.completePath(f.s.savePath)+f.completePath(f.s.imageDir)+dir, 0755)

	if err != nil {

		f.s.notice.Error(err.Error())

		return "", err

	}

	ex, err := tools.GetExtensionName(imgUrl)

	if err != nil {

		ex = "png"

		//return "", err
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

			f.s.notice.Error("图片拓展名异常:" + imgUrl)

			//获取默认图片
			if f.DefaultImg != nil {

				return f.DefaultImg(f, item), errors.New("图片拓展名异常,使用默认图片")
			}

			return "", errors.New("图片拓展名异常")
		}

	}

	imgName := (If(dir == "", "", dir+"/")).(string) + uuidString + "." + ex

	prefix := ""

	if item.ImagePrefix != nil {

		prefix = item.ImagePrefix(f, imgName)

	}

	//自动添加斜杠
	if tools.SubStr(prefix, -1, -1) != "/" {

		prefix += "/"
	}

	var imgErr error

	if f.s.CustomDownloadFun != nil {

		imgErr = f.s.CustomDownloadFun(imgUrl, imgName, f, item)

	} else {

		//imgErr = f.s.client.R().Download(imgUrl, f.completePath(f.s.savePath)+f.completePath(f.s.imageDir)+imgName)

		_, imgErr = f.s.client.R().SetOutput(f.completePath(f.s.savePath) + f.completePath(f.s.imageDir) + imgName).Get(imgUrl)

	}

	if imgErr != nil {

		msg := imgErr.Error()

		f.s.notice.Error(msg)

		//获取默认图片
		if f.DefaultImg != nil {

			return f.DefaultImg(f, item), errors.New("图片下载异常,使用默认图片：" + imgErr.Error())
		}

		return "", errors.New("图片下载异常")

	}

	//图片压缩
	if item.ImageResizePercent != 0 {

		imgDeal := image.NewImage()

		imgRes, errRes := imgDeal.LoadImage(f.completePath(f.s.savePath) + f.completePath(f.s.imageDir) + imgName)

		if errRes != nil {

			//fmt.Println("图片压缩加载错误:" + errRes.Error())

			return "", errors.New("图片压缩加载错误:" + errRes.Error())

		}

		//if errRes == nil {

		ee := imgRes.ResizePercent(item.ImageResizePercent).Save(f.completePath(f.s.savePath) + f.completePath(f.s.imageDir) + imgName + "_copy.png")

		if ee != nil {

			return "", errors.New("图片压缩错误:" + ee.Error())
		}

		//}

	}

	return (If(item.ImagePrefix == nil, "", prefix)).(string) + imgName, nil

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

// GetHtml 从链接中获取html
func (f *Form) GetHtml(url string) (string, error) {

	//content, header, err := f.s.client.R().GetToContentWithHeader(f.GetHref(url))

	//content, err := f.s.client.R().GetToContent(f.GetHref(url))

	content, err := f.s.client.R().Get(f.GetHref(url))

	if err != nil {

		return "", err

	}

	html := content.String()

	//自动转码
	if f.DisableAutoCoding == false {

		html, err = f.DealCoding(html, content.Header())

		if err != nil {

			return "", err

		}

	}

	return html, nil
}
