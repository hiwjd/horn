package dispatcher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hiwjd/horn/consumer"
	"github.com/hiwjd/horn/state"
)

type Processser func(handler *Handler, body []byte) error

func getAddr2Uids(oid int, cid string, state state.State) map[string][]string {
	uids, err := state.GetUidsInChat(oid, cid)
	if err != nil {
		log.Printf("state.GetUidsInChat err: %s \r\n", err.Error())
		return nil
	}
	//uids := state.GetUidsByChatId(chatId)
	log.Printf(" -> 获取到对话[%s]中的uids[%v] \r\n", cid, uids)

	addr2uids := make(map[string][]string)

	for _, uid := range uids {
		addr, err := state.GetPushAddrByUid(oid, uid)
		if err != nil {
			log.Printf("state.GetPushAddrByUid err: %s \r\n", err.Error())
			continue
		}

		log.Printf("  --> 获取到用户[%s]的推送地址[%s] \r\n", uid, addr)
		if _, ok := addr2uids[addr]; !ok {
			addr2uids[addr] = make([]string, 0)
		}
		addr2uids[addr] = append(addr2uids[addr], uid)
	}

	return addr2uids
}

func getAddr2UidsInOrg(oid int, state state.State) map[string][]string {
	uids, err := state.GetSidsInOrg(oid)
	if err != nil {
		log.Printf("state.GetSidsInOrg err: %s \r\n", err.Error())
		return nil
	}

	log.Printf(" -> 获取到组[%d]中的uids[%v] \r\n", oid, uids)

	addr2uids := make(map[string][]string)

	for _, uid := range uids {
		addr, err := state.GetPushAddrByUid(oid, uid)
		if err != nil {
			log.Printf("state.GetPushAddrByUid err: %s \r\n", err.Error())
			continue
		}

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

	addr2uids := getAddr2Uids(v.Oid, v.Chat.Cid, handler.state)
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

	addr2uids := getAddr2Uids(v.Oid, v.Chat.Cid, handler.state)
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

	addr2uids := getAddr2Uids(v.Oid, v.Chat.Cid, handler.state)
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

	oid := v.Oid
	chatId := v.Event.Chat.Cid
	uid := v.From.Uid
	log.Printf(" -> 开始维护对话状态 version:%s oid:%d chatId:%s uid:%s uid:%s \r\n", v.Mid, oid, chatId, uid, uid)

	err = handler.state.CreateChat(oid, v.Mid, chatId, uid)
	if err != nil {
		log.Printf(" -> 创建对话失败: %s \r\n", err.Error())
		return err
	}

	err = handler.state.JoinChat(oid, v.Mid, chatId, uid)
	if err != nil {
		log.Printf(" -> 维护对话状态失败: %s \r\n", err.Error())
		return err
	}

	// 获取被邀请对话的人的推送地址
	addr2uids := make(map[string][]string)
	for _, uid := range v.Event.Uids {
		addr, err := handler.state.GetPushAddrByUid(oid, uid)
		log.Printf("  --> 获取到用户[%s]的推送地址[%s] err[%s] \r\n", uid, addr, err.Error())
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

	oid := v.Oid
	chatId := v.Event.Chat.Cid
	uid := v.From.Uid
	log.Printf(" -> 开始维护对话状态 version:%s chatId:%s uid:%s \r\n", v.Mid, chatId, uid)
	err = handler.state.JoinChat(oid, v.Mid, chatId, uid)
	if err != nil {
		log.Printf(" -> 维护对话状态失败: %s \r\n", err.Error())
		return err
	}

	// 通知对话中的其他人 From加入对话了
	addr2uids := getAddr2Uids(oid, v.Event.Chat.Cid, handler.state)
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

	oid := v.Oid

	addr2uids := getAddr2UidsInOrg(oid, handler.state)

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

	oid := v.Oid
	state := handler.state

	chatIds, err := state.GetChatIdsByUid(oid, v.Uid)
	if err != nil {
		return err
	}

	for _, chatId := range chatIds {
		state.LeaveChat(oid, v.Mid, chatId, v.Uid)
	}

	return nil
}
