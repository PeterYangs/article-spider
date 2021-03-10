package common

import (
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
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

	//contentType, _ :=code.Find("meta[http-equiv=\"content-type\"]").Attr("content")

	contentType, _ := code.Find("meta[charset]").Attr("charset")

	//转小写
	contentType = strings.ToLower(contentType)

	switch contentType {

	case "gbk":

		html = string(tools.ConvertToByte(html, "gbk", "utf8"))

	case "gb2312":

		html = string(tools.ConvertToByte(html, "gbk", "utf8"))

	}

	return html, nil
}
