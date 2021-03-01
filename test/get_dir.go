package main

import "article-spider/common"

func main() {

	path := "game[date:m-d]/[random:1-100]/[date:Y-m-d]"

	common.GetDir(path)

}
