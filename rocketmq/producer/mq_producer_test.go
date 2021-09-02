package producer

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
	"testing"
)

var (
	topic       = "chenjia_mq_topic"
	nameSrvAddr = []string{"39.105.153.230:9876"}
	brokerAddr  = "39.105.153.230:10911"
	groupName   = "chenjia"
)

func TestProducer(t *testing.T) {

	p, _ := rocketmq.NewProducer(producer.WithNsResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		producer.WithRetry(2))
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("hello rocketmq go client! " + strconv.Itoa(i)),
		}
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}

}
