package notice

import "fmt"

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

func NewInfo(content string) *message {

	return &message{types: Info, content: content}
}

type Notice struct {
	ch chan message
}

func NewNotice() *Notice {

	return &Notice{
		ch: make(chan message, 10),
	}
}

func (n *Notice) PushMessage(message message) {

	n.ch <- message

}

func (n *Notice) Service() {

	for m := range n.ch {

		fmt.Println(m.content)
	}

}
