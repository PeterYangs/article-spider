package conf

type conf struct {
	DetailMaxCoroutines int
	MaxStrLength        int
}

var Conf = conf{
	DetailMaxCoroutines: 30,
	MaxStrLength:        60,
}
