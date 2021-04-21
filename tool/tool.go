package tool

import (
	"encoding/json"
	"fmt"
	"github.com/PeterYangs/tools"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

//转json格式
func ConversionGame(data string, list map[string]string) string {

	var jsons []map[string]string

	array := tools.Explode(",", data)

	for _, v := range array {

		jsons = append(jsons, map[string]string{"img": v, "name": ""})

	}

	jsonStr, err := json.Marshal(jsons)

	if err != nil {

		fmt.Println(err)

		return ""
	}

	return string(jsonStr)

}

func getCategory(data string, list map[string]string) string {

	switch data {

	case "角色扮演":

		return "2"

	case "动作格斗":

		return "3"

	case "休闲益智":

		return "8"

	case "飞行射击":

		return "6"

	case "冒险解密":

		return "10"

	case "策略塔防":

		return "7"

	case "赛车竞速":

		return "9"

	case "棋牌卡牌":

		return "4"

	case "音乐游戏":

		return "12"

	case "模拟经营":

		return "5"

	case "体育竞技":

		return "11"

	case "二次元养成":

		return "13"

	}

	return "2"
}

func GetDescGame(data string, resList map[string]string) string {

	if data == "" {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resList["content"]))

		if err != nil {

			return ""
		}

		return tools.SubStr(doc.Text(), 0, 65)

	}

	return data
}
