**article-spider是一个用go编写的爬取文章工具。支持两种模式，常规爬取模式和浏览器自动化模式**

安装

git clone https://github.com/PeterYangs/article-spider.git

开始使用

**爬取文字(fileTypes.SingleField)**

```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.weixz.com",
		Channel:          "/gamexz/list_[PAGE]-0.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.GameList.wd1200.mt-20px > ul > li",
		ListHrefSelector: "div.GameListIcon > a",
		DetailFields: map[string]form.Field{
			"title":   {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},

		},
	}

	spider.Start(f)


}

```

Host:网站域名

Channel：列表规则，[PAGE]替换页码

Limit：最大爬取页码

PageStart：起始页码

ListSelector：列表选择器

ListHrefSelector：列表a标签选择器，相对于列表的选择器

DetailFields：详情页选择器，key为Excel表头

ListFields：  列表页元素选择器（如需要爬列表上的缩略图或者标题）

DetailMaxCoroutine:详情页最大协程数量，默认和最大值都为列表详情页长度

DisableAutoCoding：是否关闭自动转码（目前根据页面的meta将gbk转utf8）

ProxyAddress：代理地址（你懂得）

HttpHeader：http请求头部

CustomExcelHeader：是否开启自定义excel头部

<br/>

**爬取图片(fileTypes.SingleImage)**

```

package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.weixz.com",
		Channel:          "/gamexz/list_[PAGE]-0.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.GameList.wd1200.mt-20px > ul > li",
		ListHrefSelector: "div.GameListIcon > a",
		DetailFields: map[string]form.Field{
			"title":   {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},
			"image":{Types: fileTypes.SingleImage,Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoIcon > img",ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},

		},
	}

	spider.Start(f)


}
```

**爬取富文本(fileTypes.HtmlWithImage,可以将内容中的图片下载出来并替换原链接)**

```	
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.weixz.com",
		Channel:          "/gamexz/list_[PAGE]-0.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.GameList.wd1200.mt-20px > ul > li",
		ListHrefSelector: "div.GameListIcon > a",
		DetailFields: map[string]form.Field{
			"title":   {Types: fileTypes.SingleField, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},
                        "html": {Types: fileTypes.HtmlWithImage, Selector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentTexts > div.mobileGamesContentText", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
		},
	}

	spider.Start(f)


}
	
	
```
**爬多图(fileTypes.ListImages)**

```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//爬多图
	f := form.Form{

		Host:             "https://www.duote.com",
		Channel:          "/sort/50_0_wdow_0_[PAGE]_.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.box > div.main-left-box > div > div.bd > div > div.soft-info-lists > div",
		ListHrefSelector: " a",
		DetailFields: map[string]form.Field{
			"list_img": {Types: fileTypes.ListImages, Selector: ".print-box img"},
		},
		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}


```

**爬列表元素(ListFields)**

```

package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	
	f := form.Form{

		Host:             "https://www.duote.com",
		Channel:          "/sort/50_0_wdow_0_[PAGE]_.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.box > div.main-left-box > div > div.bd > div > div.soft-info-lists > div",
		ListHrefSelector: " a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "body > div.wrap.mt_5 > div > div.main-left-box > div.down-box > div.soft-name > div > h1"},
		},
		ListFields: map[string]form.Field{
			"img": {Types: fileTypes.SingleImage, Selector: "a > img"},
		},
		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}


```

**只爬列表**


```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//只爬列表
	f := form.Form{

		Host:             "https://www.duote.com",
		Channel:          "/sort/50_0_wdow_0_[PAGE]_.html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.wrap > div.box > div.main-left-box > div > div.bd > div > div.soft-info-lists > div",
		ListHrefSelector: " a",
		ListFields: map[string]form.Field{
			"img": {Types: fileTypes.SingleImage, Selector: "a > img"},
		},
		DetailMaxCoroutine: 1,
	}

	spider.Start(f)

}



```

**代理(ProxyAddress)**

```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	//只爬列表
	f := form.Form{

		Host:             "https://store.shopping.yahoo.co.jp",
		Channel:          "/sakuranokoi/5bb3a2a955a.html?page=[PAGE]#CentSrchFilter1",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "#itmlst > ul > li",
		ListHrefSelector: " div:nth-child(1) > div > div > a",
		DetailFields: map[string]form.Field{
			"title": {Types: fileTypes.SingleField, Selector: "#shpMain > div.gdColumns.gd3ColumnItem > div.gd3ColumnItem2 > div.mdItemName > p.elCatchCopy"},
			"img":   {Types: fileTypes.SingleImage, Selector: "#itmbasic > div.elMain > ul > li.elPanel.isNew > a > img"},
		},
		DetailMaxCoroutine: 2,
		ProxyAddress:       "socks5://127.0.0.1:4781",
		
	}

	spider.Start(f)

}

```

**设置http的header(HttpHeader)**

```

package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	f := form.Form{

		Host:             "https://www.doyo.cn",
		Channel:          "/game/2-1-[PAGE].html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.mobile_game_wrap.w1168.clearfix.bg > div > div > div.tab_box > div > div > ul > li",
		ListHrefSelector: " div > a:nth-child(1)",
		DetailFields: map[string]form.Field{
			"content": {Types: fileTypes.HtmlWithImage, Selector: "#hiddenDetail > div", ExcelHeader: "C"},
		},
		DetailMaxCoroutine: 5,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		
	}

	spider.Start(f)
}



```



**自定义excel表头(ExcelHeader)**

```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
)

func main() {

	
	f := form.Form{

		Host:             "https://www.doyo.cn",
		Channel:          "/game/2-1-[PAGE].html",
		Limit:            5,
		PageStart:        1,
		ListSelector:     "body > div.mobile_game_wrap.w1168.clearfix.bg > div > div > div.tab_box > div > div > ul > li",
		ListHrefSelector: " div > a:nth-child(1)",
		DetailFields: map[string]form.Field{
			"img":         {Types: fileTypes.SingleImage, Selector: " body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.img_logo > img", ExcelHeader: "A"},
			"title":       {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.info > h1", ExcelHeader: "B"},
			"content":     {Types: fileTypes.HtmlWithImage, Selector: "#hiddenDetail > div", ExcelHeader: "C"},
			"screenshots": {Types: fileTypes.ListImages, Selector: "#slider3 > ul img", ExcelHeader: "D"},
			"size":        {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.detail_info > div.info.clearfix > span:nth-child(1) > em", ExcelHeader: "E"},
		},
		DetailMaxCoroutine: 5,
		CustomExcelHeader:  true,
	}

	spider.Start(f)
}


```

**自定义格式转换(ConversionFormatFunc)**

```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {

	f := form.Form{

		Host:             "http://www.gj078.cn",
		Channel:          "/sports/index_[PAGE].html",
		Limit:            1,
		PageStart:        1,
		ListSelector:     "#recent-content > div",
		ListHrefSelector: " div > a",
		DetailFields: map[string]form.Field{

			"title": {Types: fileTypes.SingleField, Selector: "#main > article > header > h1", ExcelHeader: "G"},
			"content": {Types: fileTypes.HtmlWithImage, Selector: "#main > article > div.entry-content", ExcelHeader: "E", ImagePrefix: "/api/uploads", ImageDir: "news/[random:1-100]"},
			"desc":    {Types: fileTypes.Attr, Selector: "meta[name=\"description\"]", AttrKey: "content", ExcelHeader: "H", ConversionFormatFunc: getDesc},
			"keyword": {Types: fileTypes.Attr, Selector: "meta[name=\"keywords\"]", AttrKey: "content", ExcelHeader: "K"},
		},
		ListFields: map[string]form.Field{
			"img": {Types: fileTypes.SingleImage, Selector: " div > a > div > img", ExcelHeader: "F", ImageDir: "news/[random:1-100]"},
		},
		HttpHeader:        map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},
		CustomExcelHeader: true,
	}

	spider.Start(f)

}

func getDesc(data string, resList map[string]string) string {

	if data == "" {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resList["content"]))

		if err != nil {

			return ""
		}

		return doc.Text()

	}

	return data
}



```

**根据某个单字段命名图片文件夹(\[singleField:title\])**


```
package main

import (
	"article-spider/fileTypes"
	"article-spider/form"
	"article-spider/spider"
	"encoding/json"
	"fmt"
	"github.com/PeterYangs/tools"
)

func main() {

	f := form.Form{

		Host:             "https://www.doyo.cn",
		Channel:          "/game/2-1-[PAGE].html",
		Limit:            1,
		PageStart:        1,
		ListSelector:     "body > div.mobile_game_wrap.w1168.clearfix.bg > div > div > div.tab_box > div > div > ul > li",
		ListHrefSelector: " div > a:nth-child(1)",
		DetailFields: map[string]form.Field{
		
			"title": {Types: fileTypes.SingleField, Selector: "body > div.game_wrap.w1200.clearfix > div.game_l > div.game_info > div.info > h1"},	
			"screenshots": {Types: fileTypes.ListImages, Selector: "#slider3 > ul img", ExcelHeader: "D", ImageDir: "[singleField:title]"},
		
		},
		DetailMaxCoroutine: 5,
		HttpHeader:         map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"},

	}

	spider.Start(f)
}


```

**浏览器自动化模式爬取(实验中)**

```
package main

import (
	"article-spider/chromedpSpider"
	"article-spider/fileTypes"
	"article-spider/form"
)

func main() {

	f := form.Form{

		Host:                "https://www.522gg.com",
		Channel:             "/game",
		Limit:               1,
		WaitForListSelector: "body > div:nth-child(5) > div > div.row.fn_mgsx10 > div",
		ListPath:            "/html/body/div[5]/div/div[2]/div",
		ListClickPath:       "/div/div/a",
		DetailFields:        map[string]form.Field{"title": {Types: fileTypes.SingleField, Selector: "body > div:nth-child(5) > div > div > div.col-xs12.col-sm12.col-md8.col-lg8 > div:nth-child(1) > div > div > div.info.w160 > div.l > h1"}},
		NextSelector:        "body > div:nth-child(5) > div > div:nth-child(3) > div > ul > li:nth-child(13) > a",
	}

	chromedpSpider.Start(f)

}


```


**web面板**

打开dist下的exe文件运行，监听8089端口

![avatar](/web/static/web.png)



