package persist

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hiwjd/horn/consumer"
	"github.com/hiwjd/horn/state"
)

type Processser func(handler *Handler, body []byte) error

var sql = `
		INSERT INTO messages
			(mid,oid,type,cid,from_uid,from_name,from_role,text,src,width,height,size,name,event)
		VALUES
			(?,  ?,  ?,   ?,      ?,       ?,        ?,        ?,   ?,  ?,    ?,     ?,   ?,   ?)
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

	_, err = db.Exec(sql, v.Mid, v.Oid, v.Type, v.Cid, v.From.Uid, v.From.Name, v.From.Role, v.Text, "", 0, 0, 0, "", "")
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

	_, err = db.Exec(sql, v.Mid, v.Oid, v.Type, v.Cid, v.From.Uid, v.From.Name, v.From.Role, "", v.Image.Src, v.Image.Width, v.Image.Height, v.Image.Size, "", "")
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

	_, err = db.Exec(sql, v.Mid, v.Oid, v.Type, v.Cid, v.From.Uid, v.From.Name, v.From.Role, "", v.File.Src, 0, 0, v.File.Size, v.File.Name, "")
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
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

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	bs, err := json.Marshal(v.Event)
	if err != nil {
		log.Printf(" -> 把Event转成json失败: %s \r\n", err.Error())
	}

	_, err = db.Exec(sql, v.Mid, v.Oid, v.Type, v.Event.Chat.Cid, v.From.Uid, v.From.Name, v.From.Role, "", "", 0, 0, 0, "", string(bs))
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
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

	db, err := handler.mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return err
	}

	bs, err := json.Marshal(v.Event)
	if err != nil {
		log.Printf(" -> 把Event转成json失败: %s \r\n", err.Error())
	}

	_, err = db.Exec(sql, v.Mid, v.Oid, v.Type, v.Event.Cid, v.From.Uid, v.From.Name, v.From.Role, "", "", 0, 0, 0, "", string(bs))
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}

	return nil
}

func viewPageProcesser(handler *Handler, body []byte) error {
	log.Println(" -> viewPageProcesser")
	var v state.Track
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

	sql := `
		INSERT INTO tracks
			(tid, vid, fp, oid, url, title, referer, os, browser, ip, addr)
		VALUES
			(?,        ?,   ?,  ?,     ?,   ?,     ?,       ?,  ?,       ?,   ?)
	`
	r, err := db.Exec(sql, v.Tid, v.Vid, v.Fp, v.Oid, v.Url, v.Title, v.Referer, v.Os, v.Browser, v.Ip, v.Addr)
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		log.Printf(" -> 保存浏览记录出错:%s \r\n", err.Error())
		return err
	}
	log.Printf(" -> 新增访问记录，影响行数[%d] \r\n", n)

	sql = `
		INSERT INTO visitors
			(vid, oid, state, fp, tid)
		VALUES
			(?,   ?,   ?,     ?,  ?)
		ON DUPLICATE KEY UPDATE
			fp=?, tid=?, updated_at=?
	`
	r, err = db.Exec(sql, v.Vid, v.Oid, "on", v.Fp, v.Tid, v.Fp, v.Tid, time.Now())
	if err != nil {
		log.Printf(" -> 执行失败: %s \r\n", err.Error())
		return err
	}
	n, err = r.RowsAffected()
	if err != nil {
		log.Printf(" -> 保存浏览记录出错:%s \r\n", err.Error())
		return err
	}
	log.Printf(" -> 新增/更新访客信息，影响行数[%d] \r\n", n)

	return nil
}
