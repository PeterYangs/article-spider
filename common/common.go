package common

import (
	"errors"
	"github.com/PeterYangs/article-spider/fileTypes"
	"github.com/PeterYangs/article-spider/form"
	ff "github.com/PeterYangs/article-spider/form"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GetHref 获取完整a链接
func GetHref(href string, host string) string {

	case1, _ := regexp.MatchString("^/[a-zA-Z0-9_]+.*", href)

	case2, _ := regexp.MatchString("^//[a-zA-Z0-9_]+.*", href)

	case3, _ := regexp.MatchString("^(http|https).*", href)

	switch true {

	case case1:

		href = host + href

		break

	case case2:

		//获取当前网址的协议
		res := regexp.MustCompile("^(https|http).*").FindStringSubmatch(host)

		href = res[1] + ":" + href

		break

	case case3:

		break

	default:

		href = host + "/" + href
	}

	return href

}

func GetDir(path string, singleFieldMap *sync.Map) string {

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

				field_, ok := singleFieldMap.Load(v[1])

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

// DealCoding 解决编码问题
func DealCoding(html string, header http.Header) (string, error) {

	headerContentType_ := header["Content-Type"]

	if len(headerContentType_) > 0 {

		headerContentType := headerContentType_[0]

		charset := GetCharsetByContentType(headerContentType)

		switch charset {

		case "gbk":

			return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

		case "gb2312":

			return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

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

		html = string(tools.ConvertToByte(html, "gbk", "utf8"))

	case "gb2312":

		html = string(tools.ConvertToByte(html, "gbk", "utf8"))

	}

	contentType, _ = code.Find("meta[http-equiv=\"Content-Type\"]").Attr("content")

	charset := GetCharsetByContentType(contentType)

	switch charset {

	case "gbk":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	case "gb2312":

		return string(tools.ConvertToByte(html, "gbk", "utf8")), nil

	}

	return html, nil
}

// GetCharsetByContentType 从contentType中获取编码
func GetCharsetByContentType(contentType string) string {

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

// ResolveSelector 解析选择器
func ResolveSelector(form form.Form, doc *goquery.Document, selector map[string]form.Field) map[string]string {

	var res = make(map[string]string)

	var lock = sync.Mutex{}

	var wait = sync.WaitGroup{}

	var singleFieldMap = sync.Map{}

	//解析详情页面选择器
	for field, item := range selector {

		switch item.Types {

		//单个文字字段
		case fileTypes.SingleField:

			v := doc.Find(item.Selector).Text()

			//singleFieldChan <- v
			singleFieldMap.Store(field, v)

			lock.Lock()
			res[field] = v
			lock.Unlock()
			break

		//单个文字字段
		case fileTypes.Attr:

			v, _ := doc.Find(item.Selector).Attr(item.AttrKey)

			lock.Lock()
			res[field] = v
			lock.Unlock()

			break

		//只爬html（不包括图片）
		case fileTypes.OnlyHtml:

			v, err := doc.Find(item.Selector).Html()

			if err != nil {

				ErrorLine(form, err.Error())

				res[field] = ""

				break

			}

			lock.Lock()
			res[field] = v
			lock.Unlock()

			break

		//爬取html，包括图片
		case fileTypes.HtmlWithImage:

			wait.Add(1)

			go func(doc *goquery.Document, form ff.Form, item ff.Field, lock *sync.Mutex, wait *sync.WaitGroup, res *map[string]string, field string, singleFieldMap *sync.Map) {

				defer wait.Done()

				html, err := doc.Find(item.Selector).Html()

				if err != nil {

					ErrorLine(form, err.Error())

					return

				}

				htmlImg, err := goquery.NewDocumentFromReader(strings.NewReader(html))

				if err != nil {

					ErrorLine(form, err.Error())

					return

				}

				var waitImg sync.WaitGroup

				var imgList = sync.Map{}

				htmlImg.Find("img").Each(func(i int, selection *goquery.Selection) {

					//img, imgBool := selection.Attr("src")
					//
					//if imgBool == false {
					//
					//	//懒加载
					//	img, imgBool = doc.Find(item.Selector).Attr("data-original")
					//
					//	if imgBool == false {
					//
					//		ErrorLine(form, "HtmlWithImage图片选择器未找到")
					//
					//		return
					//	}
					//
					//}

					img, err := getImageLink(selection)

					if err != nil {

						ErrorLine(form, "SingleImage图片选择器未找到")

						return
					}

					waitImg.Add(1)

					go func(waitImg *sync.WaitGroup, imgList *sync.Map, singleFieldMap *sync.Map) {

						defer waitImg.Done()

						imgName := DownImg(form, img, item, singleFieldMap)

						imgList.Store(imgName, img)

					}(&waitImg, &imgList, singleFieldMap)

				})

				waitImg.Wait()

				//html = strings.Replace(html, img, imgName, 1)

				imgList.Range(func(key, value interface{}) bool {

					html = strings.Replace(html, value.(string), key.(string), -1)

					return true
				})

				lock.Lock()
				resTemp := *res
				resTemp[field] = html
				lock.Unlock()

				//resChan <- map[string]string{field:html}

			}(doc, form, item, &lock, &wait, &res, field, &singleFieldMap)

		//单个图片
		case fileTypes.SingleImage:

			wait.Add(1)

			go func(doc *goquery.Document, form ff.Form, item ff.Field, lock *sync.Mutex, wait *sync.WaitGroup, res *map[string]string, field string, singleFieldMap *sync.Map) {

				defer wait.Done()

				//imgUrl, imgBool := doc.Find(item.Selector).Attr("src")

				//if imgBool == false {
				//
				//	//懒加载
				//	imgUrl, imgBool = doc.Find(item.Selector).Attr("data-original")
				//
				//	if imgBool == false {
				//
				//		ErrorLine(form, "SingleImage图片选择器未找到")
				//
				//		return
				//	}
				//
				//}

				imgUrl, err := getImageLink(doc.Find(item.Selector))

				if err != nil {

					ErrorLine(form, "SingleImage图片选择器未找到")

					return
				}

				imgName := DownImg(form, imgUrl, item, singleFieldMap)

				lock.Lock()
				resTemp := *res
				resTemp[field] = imgName
				lock.Unlock()

			}(doc, form, item, &lock, &wait, &res, field, &singleFieldMap)

			break

		//多个图片
		case fileTypes.ListImages:

			wait.Add(1)

			go func(doc *goquery.Document, form ff.Form, item ff.Field, lock *sync.Mutex, wait *sync.WaitGroup, res *map[string]string, field string, singleFieldMap *sync.Map) {

				defer wait.Done()

				var waitImg sync.WaitGroup

				var imgList = sync.Map{}

				//sync.

				doc.Find(item.Selector).Each(func(i int, selection *goquery.Selection) {

					//imgUrl, imgBool := selection.Attr("src")
					//
					//if imgBool == false {
					//
					//	//懒加载
					//	imgUrl, imgBool = doc.Find(item.Selector).Attr("data-original")
					//
					//	if imgBool == false {
					//
					//		ErrorLine(form, "ListImages图片选择器未找到")
					//
					//		return
					//	}
					//
					//}

					imgUrl, err := getImageLink(selection)

					if err != nil {

						ErrorLine(form, "ListImages图片选择器未找到")

						return
					}

					waitImg.Add(1)

					go func(waitImg *sync.WaitGroup, imgList *sync.Map, singleFieldMap *sync.Map) {

						defer waitImg.Done()

						imgName := DownImg(form, imgUrl, item, singleFieldMap)

						imgList.Store(imgName, "")

					}(&waitImg, &imgList, singleFieldMap)

				})

				waitImg.Wait()

				var strArray []string

				imgList.Range(func(key, value interface{}) bool {

					strArray = append(strArray, key.(string))

					return true
				})

				array := tools.Join(",", strArray)

				lock.Lock()
				resTemp := *res
				resTemp[field] = array
				lock.Unlock()

			}(doc, form, item, &lock, &wait, &res, field, &singleFieldMap)

		//固定数据
		case fileTypes.Fixed:

			lock.Lock()
			res[field] = item.Selector
			lock.Unlock()

		}

	}

	wait.Wait()

	return res

}

// DownImg 下载图片（包括生产文件夹）
func DownImg(form form.Form, url string, item form.Field, singleFieldMap *sync.Map) string {

	//获取完整链接
	imgUrl := GetHref(url, form.Host)

	//生成随机名称
	uuidString := uuid.NewV4().String()

	uuidString = strings.Replace(uuidString, "-", "", -1)

	dir := ""

	if item.ImageDir != "" {

		//获取图片文件夹
		dir = GetDir(item.ImageDir, singleFieldMap)

		//panic(dir)

		//设置文件夹
		err := tools.MkDirDepth("image/" + dir)

		if err != nil {

			ErrorLine(form, err.Error())

			//return
		}
	}

	ex, err := tools.GetExtensionName(imgUrl)

	if err != nil {

		ex = "png"
	}

	allowImage := []string{"png", "jpg", "jpeg", "gif", "jfif"}

	//自定义允许下载的图片拓展名
	if len(form.AllowImageExtension) > 0 {

		allowImage = form.AllowImageExtension
	}

	if !tools.In_array(allowImage, strings.ToLower(ex)) {

		ErrorLine(form, "图片拓展名异常:"+imgUrl)

		//获取默认图片
		if item.DefaultImg != nil {

			return item.DefaultImg(form, item)
		}

		return ""
	}

	imgName := (If(dir == "", "", dir+"/")).(string) + uuidString + "." + ex

	//panic(imgName)

	imgErr := tools.DownloadFile(imgUrl, "image/"+imgName, form.HttpSetting)

	if imgErr != nil {

		//log.Println(imgErr)
		msg := imgErr.Error()

		ErrorLine(form, msg)

		//fmt.Println(imgUrl)

		//获取默认图片
		if item.DefaultImg != nil {

			return item.DefaultImg(form, item)
		}

		return ""

	}

	return (If(item.ImagePrefix == "", "", item.ImagePrefix+"/")).(string) + imgName

}

// ResolveFields 解析字段
func ResolveFields(field map[string]interface{}) map[string]form.Field {

	fields := make(map[string]form.Field)

	for i, v := range field {

		item := v.(map[string]interface{})

		types := item["types"]

		fields[i] = form.Field{
			Types:       fileTypes.FieldTypes((types).(float64)),
			Selector:    (item["selector"]).(string),
			ImagePrefix: item["imagePrefix"].(string),
			ImageDir:    item["imageDir"].(string),
		}
	}

	return fields

}

// ErrorLine 错误日志
func ErrorLine(form form.Form, msg string) {

	_, f, l, _ := runtime.Caller(1)

	fullMsg := msg + " in " + f + strconv.Itoa(l)

	//fmt.Println("输出")

	form.BroadcastChan <- map[string]string{"types": "error", "data": fullMsg}

}

// ConversionFormat 处理格式转换
func ConversionFormat(form ff.Form, resList map[string]string) map[string]string {

	tempRes := resList

	//合并列表和详情选择器
	var all = make(map[string]ff.Field)

	for i, v := range form.ListFields {

		all[i] = v

	}

	for i, v := range form.DetailFields {

		all[i] = v

	}

	for i, v := range all {

		if v.ConversionFormatFunc != nil {

			tempRes[i] = v.ConversionFormatFunc(resList[i], resList)

		}

	}

	return tempRes

}

// GetChannelList 获取栏目链接
func GetChannelList(form form.Form, callback func(listUrl string)) {

	if form.ChannelFunc == nil {

		//当前页码
		var pageCurrent int

		form.Progress.Store("maxPage", float32(form.Limit-form.PageStart+1))
		form.Progress.Store("currentPage", float32(0))

		for pageCurrent = form.PageStart; pageCurrent <= form.Limit; pageCurrent++ {

			//当前列表url
			url := form.Host + strings.Replace(form.Channel, "[PAGE]", strconv.Itoa(pageCurrent), -1)

			callback(url)

			currentPage, _ := form.Progress.Load("currentPage")

			//这里有点恶心，有没有简单的写法
			c := currentPage.(float32)
			c++
			form.Progress.Store("currentPage", c)

		}

		return
	}

	//自定义栏目
	for _, i := range form.ChannelFunc(form) {

		callback(form.Host + i)

	}

}

//获取图片链接
func getImageLink(imageDoc *goquery.Selection) (string, error) {

	imgUrl, imgBool := imageDoc.Attr("src")

	if imgBool == false {

		//懒加载
		imgUrl, imgBool = imageDoc.Attr("data-original")

		if imgBool == false {

			//ErrorLine(form, "ListImages图片选择器未找到")

			return "", errors.New("未找到图片链接")
		}

	}

	return imgUrl, nil
}

// OnlyList 只爬列表
func OnlyList(form ff.Form, s *goquery.Selection) bool {

	if len(form.DetailFields) <= 0 && len(form.ListFields) > 0 {

		ts, err := s.Html()

		if err != nil {

			ErrorLine(form, err.Error())

			return true

		}

		tempDoc, err := goquery.NewDocumentFromReader(strings.NewReader(ts))

		if err != nil {

			ErrorLine(form, err.Error())

			return true
		}

		res := ResolveSelector(form, tempDoc, form.ListFields)

		form.Storage <- res

		return true

	}

	return false
}
