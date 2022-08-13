package article_spider

type Mode int

const (
	Normal Mode = 0x00000 //常规模式
	Api    Mode = 0x00001 //api模式
	Auto   Mode = 0x00002 //自动化模式
	Url    Mode = 0x00003 //详情页链接模式
)

func (m Mode) ToString() string {

	switch m {

	case 0x00000:

		return "Normal"

	case 0x00001:

		return "Api"

	case 0x00002:

		return "Auto"

	case 0x00003:

		return "Url"

	}

	return ""
}

type NextPageMode int

const (
	Pagination NextPageMode = 0 //常规分页
	LoadMore   NextPageMode = 1 //加载更多
)
