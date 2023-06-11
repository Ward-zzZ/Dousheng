package mq

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tiktok-demo/cmd/favorite/config"
	"tiktok-demo/cmd/favorite/pkg/mysql"
	"tiktok-demo/shared/consts"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/streadway/amqp"
)

// TODO: refactor this package to make it more generic and reusable.
var MysqlManager *mysql.FavoriteManager
var FavoriteActor *Actor

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
	FavoriteActor, err = NewActor(amqpConn, "favorite_action")
	if err != nil {
		klog.Fatal("cannot create add actor", err)
	}
	// start consumer thread
	go FavoriteActor.Consumer(context.Background())

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
	go a.FavoriteAction(msgs)

	<-forever
}

// the consumer function
func (a *Actor) FavoriteAction(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		klog.Infof("Received a message: %s", d.Body)
		// parse message
		message := string(d.Body)
		messageSlice := strings.Split(message, ":")
		userId, err := strconv.ParseInt(messageSlice[0], 10, 64)
		if err != nil {
			klog.Errorf("message format error: %s", message)
			continue
		}
		videoId, err := strconv.ParseInt(messageSlice[1], 10, 64)
		if err != nil {
			klog.Errorf("message format error: %s", message)
			continue
		}
		actionType, err := strconv.Atoi(messageSlice[2])
		if err != nil {
			klog.Errorf("message format error: %s", message)
		}
		klog.Infof("Mq update favorite: %d, %d, %d", userId, videoId, actionType)

		// update mysql
		if err := MysqlManager.FavoriteAction(userId, videoId, actionType == 1); err != nil {
			klog.Errorf("update favorite error: %v", err)
		}
	}
}
