package conf

type conf struct {
	DetailMaxCoroutines int
}

var Conf = conf{
	DetailMaxCoroutines: 30,
}
