package admin

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"testing"
)

var (
	topic       = "chenjia_mq_topic"
	nameSrvAddr = []string{"39.105.153.230:9876"}
	brokerAddr  = "39.105.153.230:10911"
)

func TestTopic(t *testing.T) {
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)))
	if err != nil {
		fmt.Println(err.Error())
	}
	createTopic(testAdmin, topic)
	//deleteTopic(testAdmin)
	//shutdownAdmin(testAdmin)
}

//create topic
func createTopic(opt admin.Admin, topic string) {
	err := opt.CreateTopic(context.Background(), admin.WithTopicCreate(topic), admin.WithBrokerAddrCreate(brokerAddr))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//delete topic
func deleteTopic(opt admin.Admin) {
	err := opt.DeleteTopic(context.Background(), admin.WithTopicDelete(topic), admin.WithBrokerAddrDelete(brokerAddr), admin.WithNameSrvAddr(nameSrvAddr))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//close admin
func shutdownAdmin(opt admin.Admin) {
	err := opt.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
}
