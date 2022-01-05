package mode

type Mode int

const (
	Normal Mode = 0x00000 //常规模式
	Api    Mode = 0x00001 //api模式
	Auto   Mode = 0x00002 //自动化模式
)

type NextPageMode int

const (
	Pagination NextPageMode = 0 //常规分页
	LoadMore   NextPageMode = 1 //加载更多
)
