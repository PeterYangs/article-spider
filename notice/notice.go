package notice

import (
	"fmt"
)

type types int

const (
	Info  = 0x00000
	Debug = 0x00001
	Error = 0x00002
)

type message struct {
	types   types
	content interface{}
}

func NewInfo(content interface{}) *message {

	return &message{types: Info, content: content}
}

func NewError(content interface{}) *message {

	return &message{types: Error, content: content}
}

func NewDebug(content interface{}) *message {

	return &message{types: Debug, content: content}
}

type Notice struct {
	ch chan *message
}

func NewNotice() *Notice {

	ch := make(chan *message, 10)

	return &Notice{
		ch: ch,
	}
}

func (n *Notice) PushMessage(message *message) {

	n.ch <- message

}

func (n *Notice) Service(closeEvent func()) {

	defer func() {

		closeEvent()
	}()

	for m := range n.ch {

		fmt.Println(m.content)
	}

}
func (n *Notice) Close() {

	close(n.ch)

}
