package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hiwjd/horn/consumer/persist"
	"github.com/hiwjd/horn/mysql"

	"github.com/BurntSushi/toml"
	"github.com/nsqio/go-nsq"
)

type Config struct {
	Channel          string
	Topic            string
	NsqdTCPAddrs     []string
	LookupdHTTPAddrs []string
	MaxInFlight      int
	MysqlConfigs     map[string]*mysql.Config
}

var (
	configPath string
	config     Config
)

func init() {
	flag.StringVar(&configPath, "c", "./persist.toml", "配置文件的路径")
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

	mysqlManager := mysql.New(config.MysqlConfigs)

	consumer.AddConcurrentHandlers(persist.NewHandler(mysqlManager), 4)

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
