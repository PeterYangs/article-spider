package notice

import (
	"fmt"
	"runtime/debug"
)

type types int

const (
	Info    = 0x00000
	Debug   = 0x00001
	Error   = 0x00002
	Log     = 0x00003
	Process = 0x00004
)

type message struct {
	types   types
	content []interface{}
}

func (n *Notice) Info(content ...interface{}) {

	n.ch <- &message{types: Info, content: content}
}

func (n *Notice) Error(content ...interface{}) {

	if n.debug {

		content = append(content, string(debug.Stack()))

	}

	n.ch <- &message{types: Error, content: content}
}

func (n *Notice) Debug(content ...interface{}) {

	n.ch <- &message{types: Debug, content: content}
}

func (n *Notice) Log(content ...interface{}) {

	//fmt.Println("nice啊")

	n.ch <- &message{types: Log, content: content}
}

func (n *Notice) Process(content ...interface{}) {

	n.ch <- &message{types: Process, content: content}
}

type Notice struct {
	ch chan *message
	//spider *spider.Spider
	//form *form.Form
	debug bool
}

func NewNotice() *Notice {

	ch := make(chan *message, 10)

	return &Notice{
		ch: ch,
		//spider: s,

	}
}

func (n *Notice) SetDebug(debug bool) {

	n.debug = debug
}

//func (n *Notice) PushMessage(message *message) {
//
//	n.ch <- message
//
//}

func (n *Notice) Service(closeEvent func()) {

	defer func() {

		closeEvent()
	}()

	for m := range n.ch {

		//fmt.Println(m.content...)

		//_ = m

		switch m.types {
		case Log:
			fmt.Print("\033[u\033[K")
			fmt.Println(m.content...)

		case Debug:
			fmt.Print("\033[u\033[K")
			fmt.Println(m.content...)

		case Process:
			fmt.Print("\033[u\033[K")
			fmt.Print(m.content...)
			fmt.Print("\r")

		default:
			fmt.Print("\033[u\033[K")
			fmt.Println(m.content...)
			fmt.Println()
		}

	}

}
func (n *Notice) Close() {

	close(n.ch)

}
