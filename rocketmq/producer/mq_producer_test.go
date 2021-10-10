package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var (
	topic = "qiyee-job-msg-push"
	// nameSrvAddr = []string{"172.16.5.37:9876", "172.16.5.38:9876"}
	nameSrvAddr = []string{"172.16.5.45:9876", "172.16.5.46:9876"}
)

func TestProducer(t *testing.T) {

	p, _ := rocketmq.NewProducer(producer.WithNsResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		producer.WithRetry(2))
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	// for i := 0; i < 10; i++ {
	msg := &primitive.Message{
		Topic: topic,
		Body:  []byte(builderMsg()),
	}
	res, err := p.SendSync(context.Background(), msg)

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	} else {
		fmt.Printf("send message success: result=%s\n", res.String())
	}
	fmt.Printf("生产消息第一 %d 条 \n", 1)
	time.Sleep(time.Second * time.Duration(2))
	// }
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}

}

func builderMsg() string {
	msg := Msg{
		Topic:         "qiyee-job-msg-push",
		Title:         "inbox test",
		Content:       "chenjia inbox test with roketmq",
		Summary:       "mq test",
		FromUserId:    10086,
		ToUserId:      10010,
		MsgType:       26005,
		KeyPair:       map[string]string{},
		PushChanneles: []string{"inbox"},
	}

	jsons, errs := json.Marshal(msg) //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println(string(jsons))
	return string(jsons)
}
