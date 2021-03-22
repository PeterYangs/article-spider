package main

import (
	"fmt"
	"sync"
)

func main() {

	var connectList = sync.Map{}

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	go func() {

		connectList.Store("gg", "oo")
	}()

	//time.Sleep(1*time.Second)

	v, ok := connectList.Load("xx")

	if ok {

		fmt.Println(v)
		fmt.Println("-----------")
	}

}
