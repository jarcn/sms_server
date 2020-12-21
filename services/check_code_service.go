package services

import (
	"errors"
	"fmt"
	yp "github.com/yunpian/yunpian-go-sdk/sdk"
	"log"
	"sms_server/conf"
	"sms_server/models"
	"sms_server/utils"
	"strings"
	"time"
)

var expiredMap = utils.NewExpiredMap()

func SendCode(phoneNum string, codeType int) (bool, error) {
	client := yp.New(conf.Cfg.YPAppKey)
	defer client.Close()
	param := yp.NewParam(2)
	param[yp.MOBILE] = phoneNum
	checkCode := utils.GetRandomNum(6)
	if msg, err := smsBuilder(checkCode); err != nil {
		log.Printf("短信发送失败:%s", msg)
		return false, err
	} else {
		param[yp.TEXT] = msg
	}
	log.Printf("云片请求参数:%s", param)
	r := client.Sms().SingleSend(param)
	log.Printf("云片响应参数:%s", r)
	if r.Code != 0 {
		return false, errors.New(r.Msg)
	}
	expiredMap.Set(phoneNum, checkCode, conf.Cfg.KeyTimeOut)
	saveRecord(phoneNum, checkCode, codeType)
	return true, nil
}

func saveRecord(phoneNum, checkCode string, codeType int) {
	record := models.TSmsRecord{
		PhoneNum: phoneNum,
		SmsCode:  checkCode,
		SendTime: time.Now().Unix(),
		Type:     codeType,
	}
	record.Insert()
}

func CheckCode(phoneNum, checkCode string) (bool, error) {
	if exist, code := expiredMap.Get(phoneNum); !exist {
		return false, errors.New("验证码已失效")
	} else if code != checkCode {
		return false, errors.New("验证码错误")
	}
	return true, nil
}

func smsBuilder(checkCode string) (string, error) {
	msg := make([]string, 6)
	msg = append(msg, "【环球鑫彩】")
	msg = append(msg, "您的验证码是")
	msg = append(msg, checkCode)
	msg = append(msg, "。如非本人操作，请忽略本短信")
	return strings.Replace(strings.Trim(fmt.Sprint(msg), "[]"), " ", "", -1), nil
}
