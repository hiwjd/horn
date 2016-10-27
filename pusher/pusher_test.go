package pusher

import (
	"fmt"
	"testing"
	"time"
)

func Test_Pusher(t *testing.T) {
	pusher := New(10, 128, 10)

	uid1 := "uid1"
	err := pusher.Add(uid1)
	if err != nil {
		t.Error(err)
		return
	}

	err = pusher.Push(uid1, []byte(`{"text":"msg1"}`))
	t.Log(err)
	if err != nil && err != ErrFetchTimeout && err != ErrNoNewMsg {
		t.Error(err)
		return
	}

	bs, err := pusher.Fetch(uid1, "t1", time.Second*time.Duration(1))
	t.Log(err)
	if err != nil && err != ErrFetchTimeout && err != ErrNoNewMsg {
		t.Error(err)
		return
	}
	t.Log(string(bs))

	err = pusher.Push(uid1, []byte(`{"text":"msg2"}`))
	t.Log(err)
	err = pusher.Push(uid1, []byte(`{"text":"msg3"}`))
	t.Log(err)
	bs, err = pusher.Fetch(uid1, "t1", time.Second*time.Duration(1))
	if err != nil && err != ErrFetchTimeout && err != ErrNoNewMsg {
		t.Error(err)
		return
	}
	t.Log(string(bs))

	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-quit:
				break
			default:
				bs, err := pusher.Fetch(uid1, "t1", time.Duration(2))
				if err != nil {
					t.Log(err)
				} else {
					t.Log(string(bs))
				}
			}
		}
	}()

	i := 4
	for {
		time.Sleep(1)
		pusher.Push(uid1, []byte(fmt.Sprintf(`{"text":"msg%d"}`, i)))
		i++
		if i > 8 {
			quit <- true
			break
		}
	}
}
