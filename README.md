### article-spider是一个用go编写的爬取文章工具。支持两种模式，常规爬取模式和浏览器自动化模式

[中文文档](https://www.kancloud.cn/peter_yang/article-spiderv3/2624485)
<hr/>

声明：该爬虫仅供学习使用，如产生任何法律后果，本人概不负责

**安装**

```shell
go get github.com/PeterYangs/article-spider/v4
```

[v1版本](https://github.com/PeterYangs/article-spider/tree/v1)

[v2版本](https://github.com/PeterYangs/article-spider/tree/v2)




**快速开始**

```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.925g.com/",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       2,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "app", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "app"
			}},
			"content": {ExcelHeader: "I", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "/api"
			}},
		},
		ListFields:            map[string]articleSpider.Field{},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 5,
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal, context.Background())

	s.Start()

}
```

**常用属性**

```
	Host                       string                                   //网站域名
	Channel                    string                                   //栏目链接，页码用[PAGE]替换
	PageStart                  int                                      //页码起始页
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

```

<br>

**设置header(包含cookie)**

```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.925g.com/",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       2,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "app", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "app"
			}},
			"content": {ExcelHeader: "I", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "/api"
			}},
		},
		ListFields:            map[string]articleSpider.Field{},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 5,
		HttpHeader: map[string]string{
			"cookie":     "xx",
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal,context.Background())

	s.Start()

}

```

**自定义分页链接**

```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host: "https://www.925g.com",
		ChannelFunc: func(form *articleSpider.Form) []string {

			return []string{
				"/zixun_page1.html/",
				"/zixun_page2.html/",
				"/zixun_page3.html/",
				"/zixun_page4.html/",
				"/zixun_page5.html/",
				"/zixun_page6.html/",
				"/zixun_page7.html/",
				"/zixun_page8.html/",
				"/zixun_page9.html/",
			}
		},
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       2,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img": {ExcelHeader: "H", Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "app", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "app"
			}},
			"content": {ExcelHeader: "I", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "/api"
			}},
		},
		ListFields:            map[string]articleSpider.Field{},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 5,
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal,context.Background())

	s.Start()

}

```

**详情页中间层**

```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:           "https://www.ahjingcheng.com",
		Channel:        "/show/dongzuo--------[PAGE]---/",
		ListSelector:   "body > div:nth-child(5) > div > div.col-lg-wide-75.col-xs-1.padding-0 > div:nth-child(2) > div > div.stui-pannel_bd > ul > li",
		HrefSelector:   " div > a",
		PageStart:      1,
		Length:         2,
		MiddleSelector: []string{"body > div:nth-child(3) > div > div.col-lg-wide-75.col-xs-1.padding-0 > div:nth-child(1) > div > div:nth-child(2) > div.stui-content__thumb > a"},
		DetailFields: map[string]articleSpider.Field{
			"url": {Types: articleSpider.Regular, Selector: `"url":"([0-9A-Za-z/\\._:]+)","url_next"`, RegularIndex: 1},
		},

		DetailCoroutineNumber: 1,
		HttpHeader: map[string]string{
			"cookie":     "Hm_lvt_66246be1ec92d6574526bda37cf445cc=1633767654; Hm_lvt_56a5b64a8f7a92a018377c693e064bdf=1633767654; recente=%5B%7B%22vod_name%22%3A%22%E4%B8%80%E7%BA%A7%E6%8C%87%E6%8E%A7%22%2C%22vod_url%22%3A%22https%3A%2F%2Fwww.ahjingcheng.com%2Fplay%2F119516-1-1%2F%22%2C%22vod_part%22%3A%22%E6%AD%A3%E7%89%87%22%7D%2C%7B%22vod_name%22%3A%22%E5%85%BB%E8%80%81%E5%BA%84%E5%9B%AD%22%2C%22vod_url%22%3A%22https%3A%2F%2Fwww.ahjingcheng.com%2Fplay%2F119506-1-1%2F%22%2C%22vod_part%22%3A%221080P%22%7D%2C%7B%22vod_name%22%3A%22%E4%B8%96%E7%95%8C%E4%B8%8A%E6%9C%80%E7%BE%8E%E4%B8%BD%E7%9A%84%E6%88%91%E7%9A%84%E5%A5%B3%22%2C%22vod_url%22%3A%22https%3A%2F%2Fwww.ahjingcheng.com%2Fplay%2F59426-1-1%2F%22%2C%22vod_part%22%3A%22%E5%85%A8%E9%9B%86%22%7D%2C%7B%22vod_name%22%3A%22%E6%9C%BA%E6%A2%B0%E5%B8%882%EF%BC%9A%E5%A4%8D%E6%B4%BB%E8%8B%B1%E6%96%87%E7%89%88%22%2C%22vod_url%22%3A%22https%3A%2F%2Fwww.ahjingcheng.com%2Fplay%2F91322-1-1%2F%22%2C%22vod_part%22%3A%22%E9%AB%98%E6%B8%85%22%7D%5D; Hm_lvt_66246be1ec92d6574526bda37cf445cc=1633767654; Hm_lvt_56a5b64a8f7a92a018377c693e064bdf=1633767654; PHPSESSID=7sfu1ui3crco1a817vocccl2u1; Hm_lpvt_66246be1ec92d6574526bda37cf445cc=1633914645; Hm_lpvt_56a5b64a8f7a92a018377c693e064bdf=1633914645",
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal,context.Background())

	s.Start()

}

```

**自行处理爬取结果**

```go
package main

import (
	"fmt"
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.925g.com",
		Channel:      "/zixun_page[PAGE].html/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: " a",
		PageStart:    1,
		Length:       10,
		ListFields: map[string]articleSpider.Field{

			"title": {ExcelHeader: "K", Types: articleSpider.Text, Selector: " a > div > span"},
		},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 2,
		ResultCallback: func(item map[string]string, form *articleSpider.Form) {

			for s2, s3 := range item {

				fmt.Println(s2, ":", s3)

			}

		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal,context.Background())

	s.Start()

}

```

**爬取列表是api的网页**

```go
package main

import (
	"context"
	"encoding/json"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:      "http://www.tiyuxiu.com",
		Channel:   "/data/list_0_[PAGE].json?__t=16339338",
		PageStart: 1,
		Length:    10,
		DetailFields: map[string]articleSpider.Field{

			"title":   {Types: articleSpider.Text, Selector: "h1"},
			"content": {Types: articleSpider.HtmlWithImage, Selector: "#main-content"},
		},
		//CustomExcelHeader:     true,
		DetailCoroutineNumber: 2,
		ApiConversion: func(html string, form *articleSpider.Form) []string {

			type list struct {
				Url string
			}

			var l []list

			json.Unmarshal([]byte(html), &l)

			var temp []string

			for _, l2 := range l {

				temp = append(temp, l2.Url)

			}

			return temp

		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Api,context.Background()).Debug()

	s.Start()
}
```

**自动化模式**

```go
package main

import (
	"context"
	"fmt"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	s := articleSpider.NewSpider(articleSpider.Form{

		Host:         "https://www.925g.com",
		Channel:      "/zixun/",
		ListSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li",
		HrefSelector: "  a",
		//下一页选择器
		AutoNextSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > ul > li:nth-child(11) > a",
		//列表等待选择器
		//AutoListWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.bdDiv > div > ul > li:nth-child(1)",
		//详情等待选择器
		AutoDetailWaitSelector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1",
		Length:                 3,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "J", Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"content": {ExcelHeader: "H", Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImageDir: "app", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "app"
			}},
		},

		//cookie
		HttpHeader: map[string]string{
			"cookie": "user_cookie=Vmod7XlkHN; UM_distinctid=17b805b421c1e0-0005d3dc1ac8ea-c343365-1fa400-17b805b421dda7; url_data=https://www.925g.com/zixun/,https://www.925g.com/; PHPSESSID=3m0ee50ba4r40jq3fleob2n71i; CNZZDATA1278942394=1852940385-1600066493-%7C1635143024; Hm_lvt_46233f03c62deb1e98a07bf1e1708415=1634807167,1634887947,1634955841,1635153418; Hm_lpvt_46233f03c62deb1e98a07bf1e1708415=1635153430",
		},
	}, articleSpider.Auto,context.Background())

	err := s.Start()

	if err != nil {

		fmt.Println(err)
	}

}
```

**自动化模式爬取加载更多页面**
```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"

	"github.com/chromedp/chromedp"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.btcfans.com",
		Channel:      "/zh-cn/wallet",
		ListSelector: "body > div.page-width.page-content > div.main-content > div > div.module-content > ul > li",
		HrefSelector: " a",
		//下一页选择器
		AutoNextSelector: "body > div.page-width.page-content > div.main-content > div > div.module-content > a",
		//列表等待选择器
		AutoListWaitSelector: "body > div.page-width.page-content > div.main-content > div > div.module-content > ul > li:nth-child(1)",
		//详情等待选择器
		AutoDetailWaitSelector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.info_1 > div.name > div.name-ch",
		Length:                 4,
		DetailFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "G", Types: articleSpider.Text, Selector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.info_1 > div.name > div.name-ch"},
			"content": {Types: articleSpider.HtmlWithImage, Selector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.wallet-des > div > p", ExcelHeader: "E", ImagePrefix: func(form *articleSpider.Form, imageName string) string {

				return "/api/uploads"
			}, ImageDir: "game[date:md]/[random:1-100]"},
			"desc":    {Types: articleSpider.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H"},
			"keyword": {Types: articleSpider.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
			"img":     {Types: articleSpider.Image, Selector: "body > div.page-width.page-content > div.main-content > div.wallet-detail-page > div.info_1 > div.cover > img", ExcelHeader: "F", ImageDir: "game[date:md]/[random:1-100]"},
			"type":    {Types: articleSpider.Fixed, Selector: "2", ExcelHeader: "L"},
			//"size":    {Types: fileTypes.SingleField, Selector: "#dinfo > p.base > i:nth-child(3)", ExcelHeader: "M"},
		},

		//cookie
		HttpHeader: map[string]string{
			"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36",
			"cookie":     "lang=zh-CN; lang=zh-CN; lang=zh-CN; _ga=GA1.1.1532009431.1641283813; UM_distinctid=17e24238a22739-0fc0995e9cfdad-c343365-1fa400-17e24238a2352e; guid=cff3a072d6ca30b80ee729f0884a8596f65d9a28; CNZZDATA5291371=cnzz_eid%3D1358048227-1641278212-%26ntime%3D1641338428; CNZZDATA1278599438=848177868-1641279863-%7C1641340242; Hm_lvt_ddaa34551214df42d1e5f11974f9f744=1641283822,1641346329; _csrf=3f62bc78510faa5fecfbf404cbee0ec56d1c4f3a; s_a=1; _ga_76F07DJEB4=GS1.1.1641346328.3.1.1641346978.0; Hm_lpvt_ddaa34551214df42d1e5f11974f9f744=1641346980",
		},
		//下一页模式
		AutoNextPageMode:  articleSpider.LoadMore,
		CustomExcelHeader: true,
		//爬取前置事件
		AutoPrefixEvent: func(chromedpCtx context.Context) {

			//关闭弹窗
			chromedp.Run(
				chromedpCtx,

				chromedp.Click("#Alert > div > div.sure_btn", chromedp.ByQuery),
			)

		},
	}

	s := articleSpider.NewSpider(f, articleSpider.Auto,context.Background())

	s.Start()

}
```

**代理**

```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host:         "https://www.cgcosplay.jp",
		Channel:      "/product-list?page=[PAGE]",
		ListSelector: "#inner_main_container > section > div > div.page_contents.clearfix.alllist_contents > div > div.itemlist_box.tiled_list_box.layout_photo > div > ul > li",
		HrefSelector: " div > a",
		PageStart:    1,
		Length:       10,
		ListFields: map[string]articleSpider.Field{
			"title": {ExcelHeader: "A", Types: articleSpider.Text, Selector: "div > a > div > div.list_item_data > p.item_name > span.goods_name"},
			"price": {ExcelHeader: "B", Types: articleSpider.Text, Selector: "div > a > div > div.list_item_data > div > div > p.selling_price > span.figure"},
			"img": {ExcelHeader: "C", Types: articleSpider.Image, Selector: "  div > a > div > div.list_item_photo > div > div", ImageDir: "cgcosplay_image", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "cgcosplay_image"
			}},
		},
		CustomExcelHeader:     true,
		DetailCoroutineNumber: 10,
		LazyImageAttrName:     "data-src",
		HttpProxy:             "http://127.0.0.1:4780",
	}

	s := articleSpider.NewSpider(f, articleSpider.Normal,context.Background())

	s.Start()

}
```

**排除不需要的元素**
```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
	
)

func main() {

	f := articleSpider.Form{
		Host:         "http://www.3h3.com",
		Channel:      "/news/g_38_[PAGE].html",
		ListSelector: "body > div.main > div > div > div.col-l > ul.ul-info > li",
		HrefSelector: "  div.pic > a",
		PageStart:    2,
		Length:       1,
		DetailFields: map[string]articleSpider.Field{
			"content": {Types: articleSpider.HtmlWithImage, Selector: "body > div.main > div > div > div.col-l > div.art-body", NotSelector: []string{"body > div.main > div > div > div.col-l > div.art-body > div"}},

		},

	}

	s := articleSpider.NewSpider(f, articleSpider.Normal,context.Background())

	s.Start()

}

```

**根据详情页链接爬取**

```go
package main

import (
	"context"
	articleSpider "github.com/PeterYangs/article-spider/v4"
)

func main() {

	f := articleSpider.Form{
		Host: "https://www.925g.com/",

		DetailUrls: []string{

			"https://www.925g.com/gonglue/138499.html",
			"https://www.925g.com/gonglue/138498.html",
			"https://www.925g.com/gonglue/138497.html",
			"https://www.925g.com/gonglue/138496.html",
			"https://www.925g.com/gonglue/138495.html",
			"https://www.925g.com/gonglue/138494.html",
		},
		DetailFields: map[string]articleSpider.Field{
			"title": {Types: articleSpider.Text, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.hd > h1"},
			"img":   {Types: articleSpider.Image, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd img:nth-child(1)", ImageDir: "[date:md]/[random:1-100]"},
			"content": {Types: articleSpider.HtmlWithImage, Selector: "body > div.ny-container.uk-background-default > div.wrap > div > div.commonLeftDiv.uk-float-left > div > div.articleDiv > div.bd", ImagePrefix: func(form *articleSpider.Form, path string) string {

				return "/api"
			}, ImageDir: "[date:md]/[random:1-100]"},
		},
		DetailCoroutineNumber: 3,
		FilterError:           true,
	}

	s := articleSpider.NewSpider(f, articleSpider.Url, context.Background())

	s.Start()

}
```



**关于图片保存路径说明**

**Field**中的图片路径设置


ImageDir:图片生成路径，该路径会生成在结果中，支持动态
ImagePrefix:图片前缀路径，不会出现在结果中

全局设置

SetImageDir(path),图片保存前缀，不会出现在结果中，默认是image

SetSavePath(path),图片保存文件夹，不会出现在结果中

图片保存路径拼接顺序：savePath+imageDir(全局)+imageDir(field)+文件名
图片结果路径拼接顺序: imagePrefix+ImageDir+文件名

