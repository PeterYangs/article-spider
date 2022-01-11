package article_spider

type normal struct {
	s *Spider
}

func NewNormal(s *Spider) *normal {

	return &normal{s: s}
}

func (n normal) Start() {
	//TODO implement me
	//panic("implement me")

}

func (n normal) GetList() {
	//TODO implement me
	//panic("implement me")
}

func (n normal) GetDetail() {
	//TODO implement me
	//panic("implement me")
}
