package dispatcher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hiwjd/horn/store"
)

type Processser func(handler *Handler, body []byte) error

func getAddr2Uids(chatId string, store store.Store) map[string][]string {
	uids := store.GetUidsByChatId(chatId)
	log.Printf(" -> 获取到对话[%s]中的uids[%v] \r\n", chatId, uids)

	addr2uids := make(map[string][]string)

	for _, uid := range uids {
		addr := store.GetPushAddrByUid(uid)
		log.Printf("  --> 获取到用户[%s]的推送地址[%s] \r\n", uid, addr)
		if _, ok := addr2uids[addr]; !ok {
			addr2uids[addr] = make([]string, 0)
		}
		addr2uids[addr] = append(addr2uids[addr], uid)
	}

	return addr2uids
}

func textProcesser(handler *Handler, body []byte) error {
	log.Println(" -> textProcesser")
	var v MessageText
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}
	v.T["t1"] = int(time.Now().Unix())

	addr2uids := getAddr2Uids(v.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &Message2Pusher{"text", uids, v}
		bs, er := json.Marshal(m)
		if er != nil {
			log.Printf(" -> 推送前序列化消息失败: %s \r\n", er.Error())
		} else {
			er = push(addr, bs)
			if er != nil {
				log.Printf("  --> 推送消息失败 %s: %s \r\n", bs, er.Error())
			} else {
				log.Printf("  --> 推送消息成功 addr[%s] uids[%v] %s\r\n", bs, addr, uids)
			}
		}
	}

	return nil
}

func imageProcesser(handler *Handler, body []byte) error {
	log.Println(" -> imageProcesser")
	var v MessageImage
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	addr2uids := getAddr2Uids(v.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &Message2Pusher{"image", uids, v}
		bs, er := json.Marshal(m)
		if er != nil {
			log.Printf(" -> 推送前序列化消息失败: %s \r\n", er.Error())
		} else {
			er = push(addr, bs)
			if er != nil {
				log.Printf("  --> 推送消息失败: %s \r\n", er.Error())
			} else {
				log.Printf("  --> 推送消息成功 \r\n")
			}
		}
	}

	return nil
}

func fileProcesser(handler *Handler, body []byte) error {
	log.Println(" -> fileProcesser")
	var v MessageFile
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	addr2uids := getAddr2Uids(v.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &Message2Pusher{"file", uids, v}
		bs, er := json.Marshal(m)
		if er != nil {
			log.Printf(" -> 推送前序列化消息失败: %s \r\n", er.Error())
		} else {
			er = push(addr, bs)
			if er != nil {
				log.Printf("  --> 推送消息失败: %s \r\n", er.Error())
			} else {
				log.Printf("  --> 推送消息成功 \r\n")
			}
		}
	}

	return nil
}

func eventProcesser(handler *Handler, body []byte) error {
	log.Println(" -> eventProcesser")
	return nil
}
