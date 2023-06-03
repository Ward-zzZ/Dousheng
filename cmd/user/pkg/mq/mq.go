package mq

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tiktok-demo/cmd/user/config"
	"tiktok-demo/cmd/user/pkg/mysql"
	"tiktok-demo/shared/consts"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// TODO: refactor this package to make it more generic and reusable.
var MysqlManager *mysql.UserManager
var AddActor *Actor
var DelActor *Actor

// Actor implements an amqp Actor.
type Actor struct {
	channel   *amqp.Channel
	exchange  string
	queueName string
}

func NewActor(conn *amqp.Connection, queueName string) (*Actor, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("cannot allocate channel: %v", err)
	}
	return &Actor{
		channel:   ch,
		queueName: queueName,
	}, nil
}

func InitMq() {
	c := config.GlobalServerConfig.RabbitMqInfo
	amqpConn, err := amqp.Dial(fmt.Sprintf(consts.RabbitMqURI, c.User, c.Password, c.Host, c.Port))
	if err != nil {
		klog.Fatal("cannot dial amqp", err)
	}
	AddActor, err = NewActor(amqpConn, "cnt_add")
	if err != nil {
		klog.Fatal("cannot create add actor", err)
	}
	// start consumer thread
	go AddActor.Consumer(context.Background())
	DelActor, err = NewActor(amqpConn, "cnt_del")
	if err != nil {
		klog.Fatal("cannot create del actor", err)
	}
	// start consumer thread
	go DelActor.Consumer(context.Background())

}

func (a *Actor) Publish(_ context.Context, message string) error {
	_, err := a.channel.QueueDeclare(a.queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("cannot declare queue: %v", err)
	}
	// publish message
	return a.channel.Publish(a.exchange, a.queueName, false, false, amqp.Publishing{
		Timestamp:   time.Now(),
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

func (a *Actor) Consumer(_ context.Context) {
	_, err := a.channel.QueueDeclare(a.queueName, false, false, false, false, nil)

	if err != nil {
		klog.Errorf("cannot declare queue: %v", err)
	}

	msgs, err := a.channel.Consume(a.queueName, "", true, false, false, false, nil)
	if err != nil {
		klog.Errorf("cannot Consume queue: %v", err)
	}
	var forever chan struct{}
	switch a.queueName {
	case "cnt_add":
		go a.FollowAdd(msgs)
	case "cnt_del":
		go a.FollowDel(msgs)
	}

	<-forever
}

// the consumer function of add queue.
func (a *Actor) FollowAdd(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		params := strings.Split(string(d.Body), ":")
		UserId, err := strconv.Atoi(params[0])
		if err != nil {
			klog.Errorf("Mq FollowAdd params convert error: (%v)", err)
		}
		ToUserId, err := strconv.Atoi(params[1])
		if err != nil {
			klog.Errorf("Mq FollowAdd params convert error: (%v)", err)
		}
		klog.Infof("Mq FollowAdd (%v,%v) success", UserId, ToUserId)
		if err := MysqlManager.FollowUser(int64(UserId), int64(ToUserId)); err != nil {
			klog.Errorf("Mq FollowAdd (%v,%v) error: (%v)", UserId, ToUserId, err)
		}
	}
}

// the consumer function of del queue.
func (a *Actor) FollowDel(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		params := strings.Split(string(d.Body), ":")
		UserId, err := strconv.Atoi(params[0])
		if err != nil {
			klog.Errorf("Mq FollowDel params convert error: (%v)", err)
		}
		ToUserId, err := strconv.Atoi(params[1])
		if err != nil {
			klog.Errorf("Mq FollowDel params convert error: (%v)", err)
		}
		klog.Infof("Mq FollowDel (%v,%v) success", UserId, ToUserId)
		if err := MysqlManager.UnFollowUser(int64(UserId), int64(ToUserId)); err != nil {
			klog.Errorf("Mq FollowDel (%v,%v) error: (%v)", UserId, ToUserId, err)
		}
	}
}
