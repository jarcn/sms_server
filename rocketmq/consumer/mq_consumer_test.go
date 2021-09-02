package consumer

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"os"
	"testing"
	"time"
)

var (
	topic       = "chenjia_mq_topic"
	nameSrvAddr = []string{"39.105.153.230:9876"}
	brokerAddr  = "39.105.153.230:10911"
	groupName   = "chenjia"
)

func TestConsumer(t *testing.T) {

	c, _ := rocketmq.NewPushConsumer(consumer.WithGroupName(groupName), consumer.WithNsResolver(primitive.NewPassthroughResolver(nameSrvAddr)))

	err := c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback: %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown consumer error %s", err.Error())
	}

}
