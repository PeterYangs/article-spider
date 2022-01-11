package article_spider

type SpiderMethod interface {
	Start()
	GetList()
	GetDetail()
}
