package form

import (
	"github.com/PeterYangs/article-spider/v2/fileTypes"
	"github.com/PeterYangs/article-spider/v2/mode"
	"github.com/PeterYangs/tools"
	"github.com/PeterYangs/tools/http"
	"github.com/PuerkitoBio/goquery"
	http2 "net/http"
	"regexp"
	"strings"
)

type Form struct {
	Host              string  //网站域名
	Channel           string  //栏目链接，页码用[PAGE]替换
	PageStart         int     //页码起始页
	Length            int     //爬取页码长度
	Client            *http.C //http客户端
	ListSelector      string  //列表选择器
	HrefSelector      string  //a链接选择器，相对于列表选择器
	Mode              mode.Mode
	DisableAutoCoding bool
}

type Field struct {
	Types    fileTypes.FieldTypes
	Selector string //字段选择器
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
