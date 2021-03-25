package main

import (
	"fmt"
	"runtime"
)

func main() {
	call()
}

func call() {
	//var calldepth = 1
	//fmt.Println(runtime.Caller(calldepth))

	_, f, line, _ := runtime.Caller(1)

	fmt.Println(f, line)

}
