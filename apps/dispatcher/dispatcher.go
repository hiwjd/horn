package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hiwjd/horn/dispatcher"
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
	"github.com/hiwjd/horn/store"

	"github.com/BurntSushi/toml"
	"github.com/nsqio/go-nsq"
)

var (
	configPath string
	config     dispatcher.Config
)

func init() {
	flag.StringVar(&configPath, "c", "./dispatcher.toml", "配置文件的路径")
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

	log.Println(config)

	cfg := nsq.NewConfig()

	if config.Channel == "" {
		log.Fatal("--channel is required")
	}

	if config.Topic == "" {
		log.Fatal("--topic is required")
	}

	if len(config.NsqdTCPAddrs) == 0 && len(config.LookupdHTTPAddrs) == 0 {
		log.Fatal("--nsqd-tcp-address or --lookupd-http-address required")
	}
	if len(config.NsqdTCPAddrs) > 0 && len(config.LookupdHTTPAddrs) > 0 {
		log.Fatal("use --nsqd-tcp-address or --lookupd-http-address not both")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg.UserAgent = "useragent"
	cfg.MaxInFlight = config.MaxInFlight

	consumer, err := nsq.NewConsumer(config.Topic, config.Channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	redisManager := redis.New(config.RedisConfigs)
	mysqlManager := mysql.New(config.MysqlConfigs)
	store := store.NewDefaultStore(redisManager, mysqlManager)

	//consumer.AddHandler(&TailHandler{totalMessages: *totalMessages})
	consumer.AddConcurrentHandlers(dispatcher.NewHandler(store), 4)

	err = consumer.ConnectToNSQDs(config.NsqdTCPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.ConnectToNSQLookupds(config.LookupdHTTPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-consumer.StopChan:
			return
		case <-sigChan:
			consumer.Stop()
		}
	}
}
