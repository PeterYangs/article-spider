**article-spider是一个用go编写的爬取文章工具。**

安装

git clone https://github.com/PeterYangs/article-spider.git

开始使用

**爬取文字**

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

**爬取图片**

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

**爬取富文本(html,可以将内容中的图片下载出来并替换原链接)**

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
爬多图

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

爬列表元素

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

**代理**

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

**自定义excel表头**

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


**web面板**

打开dist下的exe文件运行，监听8089端口

![avatar](/web/static/web.png)



