package models

import (
	"log"
	_ "sms_server/conf"
)

type TSmsRecord struct {
	MsgId    int64  `json:"msg_id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	PhoneNum string `json:"phone_num" xorm:"default '' comment('电话号码')  VARCHAR(11)"`
	SendTime int64  `json:"send_time" xorm:"default '' comment('短信发送时间') BIGINT(20)"`
	Type     int    `json:"type" xorm:"comment('短信类型 0:组册,1:找回密码') int(1)"`
	SmsCode  string `json:"sms_code" xorm:"default '' comment('验证码')  varchar(20)"`
}

func (sms *TSmsRecord) Insert() (int64, error) {
	session := mysqlClt.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Println("数据库连接异常", err)
		return 0, err
	}
	if result, err := session.Insert(sms); err != nil {
		log.Println("数据插入异常", err)
		return result, err
	}
	return sms.MsgId, session.Commit()
}
