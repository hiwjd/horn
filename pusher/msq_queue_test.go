package pusher

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Msg_Queue(t *testing.T) {
	assert := assert.New(t)

	l := NewLinkList(5, 2)
	assert.Equal(5, l.size, "消息队列的长度不是传入的值")

	l.Push([]byte(`{"text":"a1"}`))
	l.Push([]byte(`{"text":"b2"}`))
	//l.Push([]byte(`{"text":"c3"}`))
	//l.Push([]byte(`{"text":"d4"}`))
	//l.Push([]byte(`{"text":"e5"}`))

	t.Log(l)

	notify := l.GetNotify("t1")

	//var quit chan bool
	i := 3

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(1))
			l.Push([]byte(fmt.Sprintf(`{"text":"f%d"}`, i)))
			i++

			if i > 13 {
				//quit <- true
				break
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(3))
			select {
			case <-notify:
				ns, err := l.Fetch("t1")
				t.Log(err)
				if err == nil {
					t.Log(ns)
					j2 := NodesToJSON(ns)
					t.Log(string(j2))
					fmt.Println(string(j2))
				}
			case <-time.After(time.Second * time.Duration(5)):
				t.Log("no msg")
			}
		}
	}()

	time.Sleep(time.Second * time.Duration(12))

	// for {
	//     select {
	//         case <-quit:
	//         break
	//     }
	// }
}

func Benchmark_Queue_Push(b *testing.B) {
	b.ResetTimer()
	l := NewLinkList(3, 2)
	for i := 0; i < b.N; i++ {
		l.Push([]byte(fmt.Sprintf(`{"text":"msg%d"}`, i)))
	}
	b.Log(b.N)
	ns, err := l.Fetch("t1")
	if err != nil {
		b.Log(err)
		return
	}
	j := NodesToJSON(ns)
	b.Log(string(j))
}
