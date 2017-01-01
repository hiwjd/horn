package dispatcher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hiwjd/horn/consumer"
	"github.com/hiwjd/horn/state"
	"github.com/hiwjd/horn/utils"
)

type Processser func(handler *Handler, body []byte) error

func getAddr2Uids(oid int, cid string, state state.State) map[string][]string {
	uids, err := state.GetUidsInChat(oid, cid)
	if err != nil {
		log.Printf("state.GetUidsInChat err: %s \r\n", err.Error())
		return nil
	}

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

	addr2uids := getAddr2Uids(v.Oid, v.Cid, handler.state)
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

	addr2uids := getAddr2Uids(v.Oid, v.Cid, handler.state)
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

	addr2uids := getAddr2Uids(v.Oid, v.Cid, handler.state)
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
	state := handler.state

	// 从队列里获取到的数据解析成请求对话的数据结构
	var v consumer.MessageEventRequestChat
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	mid := v.Mid
	oid := v.Oid
	uid := v.From.Uid
	cid := v.Event.Chat.Cid
	tid := v.Event.Chat.Tid
	vid := v.Event.Chat.Vid
	sid := v.Event.Chat.Sid

	// 补全对话数据 访客信息
	v.Event.Chat.Visitor, err = state.GetVisitor(oid, vid)
	if err != nil {
		log.Printf(" -> 补全对话数据[访客信息]出错: %s \r\n", err.Error())
		return err
	}

	// 补全对话数据 客服信息
	v.Event.Chat.Staff, err = state.GetStaff(oid, sid)
	if err != nil {
		log.Printf(" -> 补全对话数据[客服信息]出错: %s \r\n", err.Error())
		return err
	}

	// 补全对话数据 访客访问轨迹
	v.Event.Chat.Tracks, err = state.GetVisitorLastTracks(oid, vid, 5)
	if err != nil {
		log.Printf(" -> 补全对话数据[访客访问轨迹]出错: %s \r\n", err.Error())
		return err
	}

	// 创建对话
	log.Printf(" -> 开始维护对话状态 version:%s oid:%d chatId:%s uid:%s sid:%s vid:%s tid:%s \r\n", v.Mid, oid, cid, uid, uid, vid, tid)
	err = handler.state.CreateChat(oid, mid, cid, uid, sid, vid, tid)
	if err != nil {
		log.Printf(" -> 创建对话失败: %s \r\n", err.Error())
		return err
	}

	// 对话参与人
	uids := make([]string, 2)
	uids[0] = v.Event.Chat.Sid
	uids[1] = v.Event.Chat.Vid

	// 获取被邀请对话的人的推送地址
	addr2uids := make(map[string][]string)
	for _, uid := range uids {
		addr, err := handler.state.GetPushAddrByUid(oid, uid)
		if err != nil {
			log.Printf("  -> 获取用户[%s]推送地址时出错[%s] \r\n", uid, err.Error())
			continue
		}
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
		log.Printf(" -> 解析消息失败: %s %s \r\n", string(body), err.Error())
		return err
	}

	oid := v.Oid
	chatId := v.Event.Cid
	uid := v.From.Uid
	log.Printf(" -> 开始维护对话状态 version:%s chatId:%s uid:%s \r\n", v.Mid, chatId, uid)
	err = handler.state.JoinChat(oid, v.Mid, chatId, uid)
	if err != nil {
		log.Printf(" -> 维护对话状态失败: %s \r\n", err.Error())
		return err
	}

	// 通知对话中的其他人 From加入对话了
	addr2uids := getAddr2Uids(oid, v.Event.Cid, handler.state)
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

func trackProcesser(handler *Handler, body []byte) error {
	log.Println(" -> trackProcesser")
	var v state.Track
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	oid := v.Oid
	addr2uids := getAddr2UidsInOrg(oid, handler.state)
	m := struct {
		Type  string       `json:"type"`
		Track *state.Track `json:"track"`
	}{
		"track",
		&v,
	}

	for addr, uids := range addr2uids {
		m := &consumer.Message2Pusher{"track", uids, m}
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
	var vs []consumer.MessageTimeout
	err := json.Unmarshal(body, &vs)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	state := handler.state

	for _, v := range vs {
		// 访客最后消息超时，结束对话，通知对话中的所有人
		// 访客心跳超时，暂时没处理
		// 客服最后消息超时，通知客服“你很久没说话了”，通知访客“客服正忙，请稍候”
		// 客服心跳超时，暂时没处理
		role := utils.GetRole(v.Uid)

		msg := consumer.MessageText{
			Message: consumer.Message{
				Type: "text",
				T:    map[string]int{"t0": 0},
				Mid:  v.Mid,
				From: &consumer.From{
					Oid:  v.Oid,
					Uid:  "0",
					Name: "系统消息",
					Role: role,
				},
				Oid:       v.Oid,
				CreatedAt: time.Now(),
			},
			Cid:  "",
			Text: "",
		}

		switch v.Type {
		case "hb": // 心跳超时了
			break
		case "lmt": // 最后的消息超时了
			if role == "visitor" {
				chatIds, err := state.GetChatIdsByUid(v.Oid, v.Uid)
				if err != nil {
					return err
				}

				for _, chatId := range chatIds {
					msg.Cid = chatId
					msg.Text = "对话结束"
					addr2uids := getAddr2Uids(v.Oid, chatId, handler.state)
					for addr, uids := range addr2uids {
						m := &consumer.Message2Pusher{"text", uids, msg}
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

					uids, err := state.GetUidsInChat(v.Oid, chatId)
					if err != nil {
						log.Printf("清理对话，获取对话中的用户时: %s \r\n", err.Error())
					} else {
						for _, uid := range uids {
							state.LeaveChat(v.Oid, v.Mid, chatId, uid)
							log.Printf("state.LeaveChat oid[%s] mid[%s] cid[%s] uid[%s] \r\n", v.Oid, v.Mid, chatId, uid)
						}
					}
				}
			} else if role == "staff" {
				chatIds, err := state.GetChatIdsByUid(v.Oid, v.Uid)
				if err != nil {
					return err
				}

				for _, chatId := range chatIds {
					msg.Cid = chatId
					msg.Text = "访客已经等了很久了"
					addr2uids := getAddr2Uids(v.Oid, chatId, handler.state)
					for addr, uids := range addr2uids {
						m := &consumer.Message2Pusher{"text", uids, msg}
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
				}
			}
			break
		}
	}

	return nil
}
