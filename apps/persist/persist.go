package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/hiwjd/horn/consumer/persist"
	"github.com/hiwjd/horn/mysql"

	"github.com/BurntSushi/toml"
	"github.com/hiwjd/horn/sendcloud"
	"github.com/nsqio/go-nsq"
)

var (
	configPath string
	config     persist.Config
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

	if len(config.Topics) < 1 {
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

	signupTpl, err := template.New("signup").Parse(config.SignupTpl)
	if err != nil {
		log.Fatalln(err)
	}
	resetpassTpl, err := template.New("resetpass").Parse(config.ResetpassTpl)
	if err != nil {
		log.Fatalln(err)
	}
	mysqlManager := mysql.New(config.MysqlConfigs)
	emailSender := sendcloud.NewEmailSender(config.SendCloudApiUser, config.SendCloudApiKey)
	handler := persist.NewHandler(mysqlManager, emailSender, signupTpl, resetpassTpl)
	consumers := make(map[*nsq.Consumer]int, len(config.Topics))

	//var consumerStoped chan *nsq.Consumer
	consumerStoped := make(chan *nsq.Consumer)

	for _, topic := range config.Topics {
		consumer, err := nsq.NewConsumer(topic, config.Channel, cfg)
		if err != nil {
			log.Fatal(err)
		}

		consumer.AddConcurrentHandlers(handler, 4)

		err = consumer.ConnectToNSQDs(config.NsqdTCPAddrs)
		if err != nil {
			log.Fatal(err)
		}

		err = consumer.ConnectToNSQLookupds(config.LookupdHTTPAddrs)
		if err != nil {
			log.Fatal(err)
		}

		consumers[consumer] = 1

		go func(consumer *nsq.Consumer) {
			select {
			case <-consumer.StopChan:
				consumerStoped <- consumer
				return
			}
		}(consumer)
	}

	for {
		select {
		case consumer := <-consumerStoped:
			delete(consumers, consumer)
			if len(consumers) == 0 {
				return
			}
		case <-sigChan:
			for consumer, _ := range consumers {
				consumer.Stop()
			}
		}
	}
}
