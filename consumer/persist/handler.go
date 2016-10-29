package persist

import (
	"errors"
	"log"

	"github.com/hiwjd/horn/mysql"

	"github.com/nsqio/go-nsq"
)

var (
	ErrUnsupportMsg = errors.New("unsupport msg")
)

type Handler struct {
	processsers  map[string]Processser
	mysqlManager *mysql.Manager
}

func NewHandler(mysqlManager *mysql.Manager) *Handler {
	ps := make(map[string]Processser, 4)
	ps["#a"] = textProcesser
	ps["#b"] = fileProcesser
	ps["#c"] = imageProcesser
	ps["#d"] = eventProcesser
	return &Handler{
		processsers:  ps,
		mysqlManager: mysqlManager,
	}
}

func (th *Handler) HandleMessage(m *nsq.Message) error {
	log.Printf("处理消息 %s \r\n", string(m.Body))
	// 1. 解析消息，验证合法性
	// 2. 处理消息的业务逻辑
	// 3. 推送消息到各pusher

	prefix := string(m.Body[0:2])
	log.Printf(" -> prefix: %s \r\n", prefix)

	var process Processser
	var ok bool
	if process, ok = th.processsers[prefix]; !ok {
		log.Printf(" -> 找不到对应的消息处理器 %s \r\n\r\n", prefix)
		m.Finish()
		return ErrUnsupportMsg
	}

	if err := process(th, m.Body[2:]); err != nil {
		log.Printf(" -> 消息处理器返回错误: %s \r\n\r\n", err.Error())
		m.Finish()
		return err
	}

	log.Printf(" -> 消息处理完毕 \r\n\r\n")

	return nil
}
