package persist

import (
	"encoding/json"
	"log"

	"github.com/hiwjd/horn/consumer"
)

type Processser func(handler *Handler, body []byte) error

var sql = `
		INSERT INTO messages
			(mid,type,chat_id,from_uid,from_name,text,src,width,height,size,name,cmd)
		VALUES
			(?,  ?,   ?,      ?,       ?,        ?,   ?,  ?,    ?,     ?,   ?,   ?)
	`

func textProcesser(handler *Handler, body []byte) error {
	log.Println(" -> textProcesser")
	var v consumer.MessageText
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	_, err = db.Exec(sql, v.Mid, v.Type, v.Chat.Id, v.From.Id, v.From.Name, v.Text, "", 0, 0, 0, "", "")
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}

	return nil
}

func imageProcesser(handler *Handler, body []byte) error {
	log.Println(" -> imageProcesser")
	var v consumer.MessageImage
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	_, err = db.Exec(sql, v.Mid, v.Type, v.Chat.Id, v.From.Id, v.From.Name, "", v.Image.Src, v.Image.Width, v.Image.Height, v.Image.Size, "", "")
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}

	return nil
}

func fileProcesser(handler *Handler, body []byte) error {
	log.Println(" -> fileProcesser")
	var v consumer.MessageFile
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	_, err = db.Exec(sql, v.Mid, v.Type, v.Chat.Id, v.From.Id, v.From.Name, "", v.File.Src, 0, 0, v.File.Size, v.File.Name, "")
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}

	return nil
}

func requestChatProcesser(handler *Handler, body []byte) error {
	log.Println(" -> requestChatProcesser")
	var v consumer.MessageCmdRequestChat
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	bs, err := json.Marshal(v.Cmd)
	if err != nil {
		log.Printf(" -> 把Cmd转成json失败: %s \r\n", err.Error())
	}

	_, err = db.Exec(sql, v.Mid, v.Type, v.Cmd.Chat.Id, v.From.Id, v.From.Name, "", "", 0, 0, 0, "", string(bs))
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}

	return nil
}

func joinChatProcesser(handler *Handler, body []byte) error {
	log.Println(" -> joinChatProcesser")
	var v consumer.MessageCmdJoinChat
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	bs, err := json.Marshal(v.Cmd)
	if err != nil {
		log.Printf(" -> 把Cmd转成json失败: %s \r\n", err.Error())
	}

	_, err = db.Exec(sql, v.Mid, v.Type, v.Cmd.Chat.Id, v.From.Id, v.From.Name, "", "", 0, 0, 0, "", string(bs))
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}

	return nil
}
