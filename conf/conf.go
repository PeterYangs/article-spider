package conf

type Conf struct {
	DetailMaxCoroutines int
	MaxStrLength        int
	ImageDir            string
}

func NewConf() *Conf {

	return &Conf{DetailMaxCoroutines: 30, MaxStrLength: 60, ImageDir: "image"}
}

//var Conf = conf{
//	DetailMaxCoroutines: 30,
//	MaxStrLength:        60,
//	ImageDir:            "image",
//}
