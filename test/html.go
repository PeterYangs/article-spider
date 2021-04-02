package main

import (
	"fmt"
	"github.com/PeterYangs/tools"
)

func main() {

	url := "https://www.crsky.com/soft/273960.html"
	//url:="https://www.duote.com/sort/50_0_wdow_0_1_.html"

	resp, _ := tools.GetToResp(url, tools.HttpSetting{})

	fmt.Println(resp.Header.Get("Content-Type"))

	//fmt.Println(html.)

	//if tools.IsGBK([]byte(html)) {
	//
	//	fmt.Println(string( tools.ConvertToByte(html,"gbk","utf8")))
	//}else {
	//
	//	fmt.Println(html)
	//}

	//fmt.Println(html)

}
