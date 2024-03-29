package article_spider

import (
	"fmt"
	"log"
	"runtime/debug"
	"strings"
)

type types int

const (
	Info    = 0x00000
	Debug   = 0x00001
	Error   = 0x00002
	Log     = 0x00003
	Process = 0x00004
	Finish  = 0x00005
)

type message struct {
	types   types
	content []interface{}
}

func (n *Notice) Info(content ...interface{}) {

	n.ch <- &message{types: Info, content: content}
}

func (n *Notice) Error(content ...interface{}) {

	if n.s.debug {

		content = append(content, string(debug.Stack()))

	}

	n.ch <- &message{types: Error, content: content}
}

func (n *Notice) Debug(content ...interface{}) {

	n.ch <- &message{types: Debug, content: content}
}

func (n *Notice) Log(content ...interface{}) {

	n.ch <- &message{types: Log, content: content}
}

func (n *Notice) Process(content ...interface{}) {

	n.ch <- &message{types: Process, content: content}
}

func (n *Notice) Finish(content ...interface{}) {

	n.ch <- &message{types: Finish, content: content}
}

type Notice struct {
	ch chan *message
	s  *Spider
}

func NewNotice(s *Spider) *Notice {

	ch := make(chan *message, 10)

	return &Notice{
		ch: ch,
		s:  s,
	}
}

func (n *Notice) Service() {

	n.s.wait.Add(1)

	defer n.s.wait.Done()

	//n.s.form.

	for {

		select {
		case m := <-n.ch:

			switch m.types {

			case Process:

				if n.s.form.DisableMessage {

					break
				}

				fmt.Print(m.content...)
				fmt.Print(strings.Repeat(" ", 50))
				fmt.Print("\r")

			case Finish:

				fmt.Println()
				fmt.Println()
				log.Println(m.content...)
				fmt.Println()
				fmt.Println()

				return

			default:

				if n.s.form.DisableMessage {

					break
				}

				fmt.Println()
				fmt.Println()
				fmt.Println(m.content...)
				fmt.Println()
			}

		}

	}

}

func (n *Notice) Stop() {

	close(n.ch)

	//fmt.Println("退出通知管道")

}
