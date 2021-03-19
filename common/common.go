package common

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	uuid "github.com/satori/go.uuid"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//获取完整a链接
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

func GetDir(path string) string {

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

	return path

}

//伪三元运算
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

//解决编码问题
func DealCoding(html string) (string, error) {

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

	//捕获编码
	r, _ := regexp.Compile(`charset=(.*)`)

	re := r.FindAllStringSubmatch(contentType, 1)

	if len(re) > 0 {

		c := re[0][1]

		switch c {

		case "gbk":

			html = string(tools.ConvertToByte(html, "gbk", "utf8"))

		case "gb2312":

			html = string(tools.ConvertToByte(html, "gbk", "utf8"))

		}

	}

	return html, nil
}

//解析选择器
func ResolveSelector(form form.Form, doc *goquery.Document, selector map[string]form.Field) map[string]string {

	var res = make(map[string]string)

	//解析详情页面选择器
	for field, item := range selector {

		switch item.Types {

		//单个文字字段
		case fileTypes.SingleField:

			v := doc.Find(item.Selector).Text()

			//fmt.Println(v)

			res[field] = v

			break

		//只爬html（不包括图片）
		case fileTypes.OnlyHtml:

			v, err := doc.Find(item.Selector).Html()

			if err != nil {

				fmt.Println(err)

				break

			}

			res[field] = v

			break

		//爬取html，包括图片
		case fileTypes.HtmlWithImage:

			html, err := doc.Find(item.Selector).Html()

			if err != nil {

				fmt.Println(err)

				break

			}

			htmlImg, err := goquery.NewDocumentFromReader(strings.NewReader(html))

			if err != nil {

				fmt.Println(err)

				break

			}

			htmlImg.Find("img").Each(func(i int, selection *goquery.Selection) {

				img, b := selection.Attr("src")

				if b == true {

					imgName := DownImg(form, img, item)

					html = strings.Replace(html, img, imgName, 1)

				}

			})

			res[field] = html

		//单个图片
		case fileTypes.SingleImage:

			imgUrl, imgBool := doc.Find(item.Selector).Attr("src")

			if imgBool == false {

				fmt.Println("SingleImage图片选择器未找到")

				break

			}

			imgName := DownImg(form, imgUrl, item)

			res[field] = imgName

			break

		//多个图片
		case fileTypes.ListImages:

			imgList := ""

			doc.Find(item.Selector).Each(func(i int, selection *goquery.Selection) {

				imgUrl, imgBool := selection.Attr("src")

				if imgBool == false {

					fmt.Println("ListImages图片选择器未找到")

				} else {

					imgName := DownImg(form, imgUrl, item)

					//imgList=append(imgList, imgName)

					imgList += imgName + ","

				}

				//fmt.Println(imgName)

			})

			res[field] = imgList

		}

	}

	return res

}

func DownImg(form form.Form, url string, item form.Field) string {

	//获取完整链接
	imgUrl := GetHref(url, form.Host)

	//生成随机名称
	uuidString := uuid.NewV4().String()

	dir := ""

	if item.ImageDir != "" {

		//获取图片文件夹
		dir = GetDir(item.ImageDir)

		//设置文件夹
		err := tools.MkDirDepth("image/" + dir)

		if err != nil {

			log.Println(err)

			//return
		}
	}

	imgName := (If(dir == "", "", dir+"/")).(string) + uuidString + "." + tools.GetExtensionName(imgUrl)

	imgErr := tools.DownloadImage(imgUrl, "image/"+imgName, form.HttpSetting)

	if imgErr != nil {

		log.Println(imgErr)

	}

	return (If(item.ImagePrefix == "", "", item.ImagePrefix+"/")).(string) + imgName

}

//解析字段
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
