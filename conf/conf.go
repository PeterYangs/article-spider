package conf

type conf struct {
	DetailMaxCoroutines int
	MaxStrLength        int
	ImageDir            string
}

var Conf = conf{
	DetailMaxCoroutines: 30,
	MaxStrLength:        60,
	ImageDir:            "cccc",
}
