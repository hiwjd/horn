package persist

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"time"

	"github.com/hiwjd/horn/consumer"
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/sendcloud"
	"github.com/hiwjd/horn/utils"
)

// 消费者多了部署有点不便啊
// 发邮件任务就先放在persist了

func signupEmailProcesser(handler *Handler, body []byte) error {
	log.Println(" -> signupEmailProcesser")
	var v consumer.MessageSignupEmail
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
		link := fmt.Sprintf("http://app.horn.com:9092/api/signup_confirm?s=%s", token)

		template := `<p>你好！欢迎来到HORN</p><p>点击以下链接来激活您的账号</p><p><a target="_blank" href="%s">%s</p><p>HONR团队</p>`
		html := fmt.Sprintf(template, link, link)
		data = &sendcloud.EmailData{
			From:     "welcome@horn.com",
			To:       v.Email,
			Subject:  "欢迎来到HORN",
			Html:     html,
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

func genSignupToken(mysqlManager *mysql.Manager, email string) (string, error) {
	db, err := mysqlManager.Get("write")
	if err != nil {
		log.Printf(" -> 获取数据库连接失败: %s \r\n", err.Error())
		return "", err
	}

	token := utils.RandString(53)
	expires_at := time.Now().Unix()

	sql := `
        INSERT INTO signup_email
            (email,token,count,expires_at)
        VALUES
            (?,?,1,?) 
        ON DUPLICATE KEY UPDATE 
            token=?, count=count+1, expires_at=?
    `
	_, err = db.Exec(sql, email, token, expires_at, token, expires_at)
	if err != nil {
		log.Printf(" -> 新增／更新signup_email出错: %s \r\n", err.Error())
		return "", err
	}

	return token, nil
}
