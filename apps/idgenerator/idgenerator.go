package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rs/xid"
)

var (
	configPath string
	config     Config
)

type Config struct {
	Addr string
}

func init() {
	flag.StringVar(&configPath, "c", "./idgenerator.toml", "配置文件的路径")
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

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		id := xid.New()
		fmt.Fprintf(w, "%s", id.String())
	})

	log.Printf("ID生成器启动，地址: %s \r\n", config.Addr)
	http.ListenAndServe(config.Addr, nil)
}
