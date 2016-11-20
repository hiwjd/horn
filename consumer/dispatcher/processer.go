package dispatcher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hiwjd/horn/consumer"
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

func getAddr2UidsInCompany(cid string, store store.Store) map[string][]string {
	uids := store.GetStaffsByCompany(cid)
	log.Printf(" -> 获取到组[%s]中的uids[%v] \r\n", cid, uids)

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
	var v consumer.MessageText
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}
	v.T["t1"] = int(time.Now().Unix())

	addr2uids := getAddr2Uids(v.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"text", uids, v}
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
	var v consumer.MessageImage
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	addr2uids := getAddr2Uids(v.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"image", uids, v}
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
	var v consumer.MessageFile
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	addr2uids := getAddr2Uids(v.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"file", uids, v}
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

func requestChatProcesser(handler *Handler, body []byte) error {
	log.Println(" -> requestChatProcesser")
	var v consumer.MessageEventRequestChat
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	cid := v.Cid
	chatId := v.Event.Chat.Id
	var uid, staffId string
	if v.From.Role == "staff" {
		uid = ""
		staffId = v.From.Id
	} else {
		uid = v.From.Id
		staffId = ""
	}
	log.Printf(" -> 开始维护对话状态 version:%s cid:%s chatId:%s uid:%s staffId:%s \r\n", v.Mid, cid, chatId, uid, staffId)

	err = handler.store.CreateChat(chatId, cid, uid, staffId)
	if err != nil {
		log.Printf(" -> 创建对话失败: %s \r\n", err.Error())
		return err
	}

	err = handler.store.JoinChat(v.Mid, chatId, uid, v.From.Role)
	if err != nil {
		log.Printf(" -> 维护对话状态失败: %s \r\n", err.Error())
		return err
	}

	// 获取被邀请对话的人的推送地址
	addr2uids := make(map[string][]string)
	for _, uid := range v.Event.Uids {
		addr := handler.store.GetPushAddrByUid(uid)
		log.Printf("  --> 获取到用户[%s]的推送地址[%s] \r\n", uid, addr)
		if _, ok := addr2uids[addr]; !ok {
			addr2uids[addr] = make([]string, 0)
		}
		addr2uids[addr] = append(addr2uids[addr], uid)
	}

	// 通知被邀请人，有人请求对话
	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"request_chat", uids, v}
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

func joinChatProcesser(handler *Handler, body []byte) error {
	log.Println(" -> joinChatProcesser")
	var v consumer.MessageEventJoinChat
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	chatId := v.Event.Chat.Id
	uid := v.From.Id
	log.Printf(" -> 开始维护对话状态 version:%s chatId:%s uid:%s \r\n", v.Mid, chatId, uid)
	err = handler.store.JoinChat(v.Mid, chatId, uid, v.From.Role)
	if err != nil {
		log.Printf(" -> 维护对话状态失败: %s \r\n", err.Error())
		return err
	}

	// 通知对话中的其他人 From加入对话了
	addr2uids := getAddr2Uids(v.Event.Chat.Id, handler.store)
	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"join_chat", uids, v}
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

func viewPageProcesser(handler *Handler, body []byte) error {
	log.Println(" -> viewPageProcesser")
	var v consumer.MessageViewPage
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	addr2uids := getAddr2UidsInCompany(v.Cid, handler.store)

	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"view_page", uids, v}
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

func timeoutProcesser(handler *Handler, body []byte) error {
	log.Println(" -> timeoutProcesser")
	var v consumer.MessageTimeout
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	store := handler.store

	chatIds := store.GetChatsByUid(v.Uid)
	for _, chatId := range chatIds {
		store.LeaveChat(v.Mid, chatId, v.Uid)
	}

	return nil
}
