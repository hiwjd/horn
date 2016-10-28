package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hiwjd/horn/pusher"

	"github.com/BurntSushi/toml"
	"golang.org/x/net/websocket"
)

var (
	configPath string
	config     Config
)

type Config struct {
	Name string
	Addr string
}

type JoinRequest struct {
	Uid string
}

type MessageRequest struct {
	Type string
	To   []string
	Data json.RawMessage
}

func init() {
	flag.StringVar(&configPath, "c", "./pusher.toml", "配置文件的路径")
}

func main() {
	flag.Parse()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("配置文件 %s 不存在, 使用默认配置 \r\n", configPath)
	} else {
		_, err := toml.DecodeFile(configPath, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	p := pusher.New(10, 128, 10)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}
		defer r.Body.Close()

		log.Printf("/push %s \r\n", bs)

		var req MessageRequest
		err = json.Unmarshal(bs, &req)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}

		for _, uid := range req.To {
			// todo 错误处理
			p.Push(uid, req.Data)
		}

		fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 0, "")
	})

	http.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}
		defer r.Body.Close()

		var req JoinRequest
		err = json.Unmarshal(bs, &req)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}

		err = p.Add(req.Uid)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}

		fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 0, "")
	})

	http.HandleFunc("/pull", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		q := r.URL.Query()

		uid := q.Get("uid")
		trackID := q.Get("track_id")
		log.Printf("/pull uid:%s track_id:%s \r\n", uid, trackID)
		if uid == "" || trackID == "" {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, "uid or trackID missing")
			return
		}

		keep := q.Get("keep")
		if keep == "" {
			keep = "15"
		}
		keepInt, err := strconv.Atoi(keep)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}

		if keepInt < 5 || keepInt > 30 {
			keepInt = 15
		}
		log.Printf(" -> keep: %d \r\n", keepInt)
		keepDuration := time.Duration(keepInt) * time.Second

		bs, err := p.Fetch(uid, trackID, keepDuration)
		if err != nil {
			if err == pusher.ErrFetchTimeout {
				fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 0, "")
			} else {
				fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			}
			return
		}

		fmt.Fprintf(w, `{"code":0,"msg":"","data":`)
		w.Write(bs)
		fmt.Fprintf(w, `}`)
	})

	ws := &PusherWSServer{p}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		handler := websocket.Handler(ws.Handle)
		handler.ServeHTTP(w, r)
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		uid := q.Get("uid")
		m := p.Stats(uid)
		bs, err := json.Marshal(m)
		if err != nil {
			fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, 1, err.Error())
			return
		}
		w.Write(bs)
	})

	log.Printf("PUSHER[%s]启动，地址: %s \r\n", config.Name, config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}

// PusherWSServer 是websocket方式的推送服务
type PusherWSServer struct {
	p *pusher.Pusher
}

// Handle 处理连接
func (c *PusherWSServer) Handle(conn *websocket.Conn) {
	log.Println("ws connected")
	q := conn.Request().URL.Query()

	uid := q.Get("uid")
	trackID := q.Get("track_id")
	if uid == "" || trackID == "" {
		fmt.Fprintf(conn, `{"code":0,"msg":"uid or track_id empty"}`)
		conn.Close()
		return
	}

	keep := q.Get("keep")
	if keep == "" {
		keep = "30"
	}
	keepInt, err := strconv.Atoi(keep)
	if err != nil {
		fmt.Fprintf(conn, `{"code":%d,"msg":"%s"}`, 1, err.Error())
		conn.Close()
		return
	}

	if keepInt < 15 || keepInt > 60 {
		keepInt = 30
	}
	keepDuration := time.Duration(keepInt) * time.Second

	for {
		bs, err := c.p.Fetch(uid, trackID, keepDuration)
		if err != nil {
			log.Println(err)
			if err == pusher.ErrFetchTimeout {
				fmt.Fprintf(conn, `{"code":0,"msg":""}`)
			} else {
				fmt.Fprintf(conn, `{"code":%d,"msg":"%s"}`, 1, err.Error())
				conn.Close()
				return
			}
		} else {
			fmt.Fprintf(conn, `{"code":%d,"msg":"","data":%s}`, 0, bs)
		}
	}
}
