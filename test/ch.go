package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan bool, 3)

	ch <- true
	ch <- true
	ch <- true
	//ch <-true

	go ccc(ch)

	for _ = range ch {

		//time.Sleep(1*time.Hour)

		//for  {

		fmt.Println("nice")
		time.Sleep(11 * time.Second)
		//}

	}

	fmt.Println("finish")

}

func ccc(ch chan bool) {

	time.Sleep(10 * time.Second)

	fmt.Println("开始关闭通道")

	close(ch)

}
