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
			"title":   {Types: fileTypes.SingleField, SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},

		},
	}

	spider.Start(f)


}

```

Host:网站域名

Channel：列表规则，[PAGE]替换页面

Limit：最大爬取页码

PageStart：起始页码

ListSelector：列表选择器

ListHrefSelector：列表a标签选择器，相对于列表的选择器

ListHrefSelector：详情页选择器，key为Excel表头

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
			"title":   {Types: fileTypes.SingleField, SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},
			"image":{Types: fileTypes.SingleImage,SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoIcon > img",ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},

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
			"title":   {Types: fileTypes.SingleField, SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentInfo.displayFlex > div.mobileGamesContentInfoText > div > h1"},
                        "html": {Types: fileTypes.HtmlWithImage, SingleSelector: "body > div.wrap > div.information-main.mt-20px.wd1200.displayFlex > div.information-main-left > div.mobileGamesContent > div.mobileGamesContentTexts > div.mobileGamesContentText", ImagePrefix: "upload", ImageDir: "[date:Ym]/[random:1-100]"},
		},
	}

	spider.Start(f)


}
	
	
```

