package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"reflect"
	"strings"
)

func main() {

	html := `
	<html>
		
		<div> <span>13</span> </div>
		<div> <span>12</span> </div>
		<div> <span>11</span> </div>
		<div> <span>10</span> </div>

</html>

`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {

		log.Fatal(err)
	}

	doc.Find("html div").Each(func(i int, selection *goquery.Selection) {

		//fmt.Println(*selection)

		//t := reflect.TypeOf(*selection)
		//v := reflect.ValueOf(*selection)
		//for k := 0; k < t.NumField(); k++ {
		//	//fmt.Printf("%s -- %v \n", t.Field(k).Name, v.Field(k))
		//
		//	fmt.Println(v.Field(k))
		//}

		//fmt.Println()

		for _, node := range selection.ToggleClass().First().Nodes {

			//fmt.Println(*node)

			t := reflect.TypeOf(*node)
			v := reflect.ValueOf(*node)
			for k := 0; k < t.NumField(); k++ {
				//fmt.Printf("%s -- %v \n", t.Field(k).Name, v.Field(k))

				fmt.Println(v.Field(k))
			}

		}
	})

	//fmt.Println(s)

}
