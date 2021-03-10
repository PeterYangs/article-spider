package spider

import (
	"article-spider/common"
	"article-spider/fileTypes"
	"article-spider/form"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"sync"
)

//爬取详情
func GetDetail(form form.Form, detailUrl string, wait *sync.WaitGroup, detailMaxChan chan int) {

	defer func(detailMaxChan chan int, max int) {

		if max != 0 {

			<-detailMaxChan

		}

		wait.Done()

	}(detailMaxChan, form.DetailMaxCoroutine)

	//获取详情页面html
	html, err := tools.GetToString(detailUrl, tools.HttpSetting{})

	//自动转码
	if form.DisableAutoCoding == false {

		html, err = common.DealCoding(html)

		if err != nil {

			fmt.Println(err)

			return

		}

	}

	//panic(html)

	if err != nil {

		fmt.Println(err)

		return

	}

	//加载
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		//log.Fatal(err)

		fmt.Println(err)

		return

	}

	var res = make(map[string]string)

	//解析详情页面选择器
	for field, item := range form.DetailFields {

		switch item.Types {

		//单个文字字段
		case fileTypes.SingleField:

			v := doc.Find(item.Selector).Text()

			fmt.Println(v)

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

	//写入管道
	form.Storage <- res

}

func DownImg(form form.Form, url string, item form.Field) string {

	//获取完整链接
	imgUrl := common.GetHref(url, form.Host)

	//生成随机名称
	uuidString := uuid.NewV4().String()

	dir := ""

	if item.ImageDir != "" {

		//获取图片文件夹
		dir = common.GetDir(item.ImageDir)

		//设置文件夹
		err := tools.MkDirDepth("image/" + dir)

		if err != nil {

			log.Println(err)

			//return
		}
	}

	imgName := (common.If(dir == "", "", dir+"/")).(string) + uuidString + "." + tools.GetExtensionName(imgUrl)

	imgErr := tools.DownloadImage(imgUrl, "image/"+imgName)

	if imgErr != nil {

		log.Println(imgErr)

	}

	return (common.If(item.ImagePrefix == "", "", item.ImagePrefix+"/")).(string) + imgName

}
