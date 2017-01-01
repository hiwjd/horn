package persist

import (
	"encoding/json"
	"errors"
	"log"

	"time"

	"bytes"

	"github.com/hiwjd/horn/consumer"
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/sendcloud"
	"github.com/hiwjd/horn/utils"
)

// 消费者多了部署有点不便啊
// 发邮件任务就先放在persist了

func sendEmailProcesser(handler *Handler, body []byte) error {
	log.Println(" -> sendEmailProcesser")
	var v consumer.MessageSendEmail
	err := json.Unmarshal(body, &v)
	if err != nil {
		log.Printf(" -> 解析消息失败: %s \r\n", err.Error())
		return err
	}

	var data *sendcloud.EmailData

	switch v.Type {
	case "signup":
		token, err := genSignupToken(handler.mysqlManager, v.Email)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		tplData := struct {
			Token string
			Email string
		}{
			Token: token,
			Email: v.Email,
		}
		err = handler.signupTpl.Execute(&buf, tplData)
		if err != nil {
			log.Printf("signup email: %s \r\n", err.Error())
		}

		data = &sendcloud.EmailData{
			From:     "team@hiyueliao.com",
			To:       v.Email,
			Subject:  "感谢您注册悦聊",
			Html:     buf.String(),
			FromName: "",
			ReplyTo:  "",
		}
	case "find_pass":
		token, err := genFindPassToken(handler.mysqlManager, v.Email)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		tplData := struct {
			Token string
			Email string
		}{
			Token: token,
			Email: v.Email,
		}
		err = handler.signupTpl.Execute(&buf, tplData)
		if err != nil {
			log.Printf("find_pass email: %s \r\n", err.Error())
		}

		data = &sendcloud.EmailData{
			From:     "team@hiyueliao.com",
			To:       v.Email,
			Subject:  "找回密码",
			Html:     buf.String(),
			FromName: "",
			ReplyTo:  "",
		}
	default:
		return errors.New("不支持的邮件发送类型")
	}

	err = handler.emailSender.SendEmail(data)
	if err != nil {
		log.Printf("邮件发送失败: %s \r\n", err.Error())
		return err
	}

	return nil
}

func genFindPassToken(mysqlManager *mysql.Manager, email string) (string, error) {
	db, err := mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return "", err
	}

	token := utils.RandString(57)
	expires_at := time.Now().Unix()

	sql := `
        INSERT INTO email_tokens
            (email,intention,token,count,state,expires_at)
        VALUES
            (?,'resetpass',?,1,?,?) 
        ON DUPLICATE KEY UPDATE 
            token=?, count=count+1, state=?, expires_at=?
    `
	_, err = db.Exec(sql, email, token, "valid", expires_at, token, "valid", expires_at)
	if err != nil {
		log.Printf(" -> 新增／更新email_tokens出错: %s \r\n", err.Error())
		return "", err
	}

	return token, nil
}

func genSignupToken(mysqlManager *mysql.Manager, email string) (string, error) {
	db, err := mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return "", err
	}

	token := utils.RandString(53)
	expires_at := time.Now().Unix()

	sql := `
        INSERT INTO email_tokens
            (email,intention,token,count,state,expires_at)
        VALUES
            (?,'signup',?,1,?,?) 
        ON DUPLICATE KEY UPDATE 
            token=?, count=count+1, state=?, expires_at=?
    `
	_, err = db.Exec(sql, email, token, "valid", expires_at, token, "valid", expires_at)
	if err != nil {
		log.Printf(" -> 新增／更新email_tokens出错: %s \r\n", err.Error())
		return "", err
	}

	return token, nil
}
