package main

import (
	articleSpider "github.com/PeterYangs/article-spider/v3"
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

	s := articleSpider.NewSpider(f, articleSpider.Normal)

	s.Start()

}
