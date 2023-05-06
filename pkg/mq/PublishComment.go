package mq

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"simple_tiktok/dal/db"
	"simple_tiktok/dal/db/model"
	"simple_tiktok/dal/redis"
	"simple_tiktok/pkg/consts"
	"simple_tiktok/pkg/errno"
)

func SendComment(comment *model.Comment) error {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial(consts.RabbitMQDSN)
	if err != nil {
		return err
	}
	// 当前函数执行完关闭连接
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// 同上，当前函数执行完就关闭连接了
	defer ch.Close()

	// 声明交换器
	err = ch.ExchangeDeclare(
		consts.CommentExchange,
		"fanout", // 还有 amqp.ExchangeTopic
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	// 将评论序列化为 JSON 格式
	commentJSON, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	// 发送消息到交换器
	err = ch.Publish(
		"",                     // 交换器名称
		consts.CommentExchange, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        commentJSON,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func ReceiveMessage(ctx context.Context) error {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial(consts.RabbitMQDSN)
	if err != nil {
		return err
	}
	// 当前函数执行完关闭连接
	defer conn.Close()

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// 同上，当前函数执行完就关闭连接了
	defer ch.Close()

	// 声明交换器
	err = ch.ExchangeDeclare(
		consts.CommentExchange,
		"fanout", // 还有 amqp.ExchangeTopic
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		consts.CommentExchange, // 队列名称
		false,                  // 持久化
		false,                  // 自动删除
		false,                  // 独占队列
		false,                  // 队列没有消费者时删除
		nil,                    // 额外参数
	)
	if err != nil {
		return errno.NewErrNo("声明消息队列失败！" + err.Error())
	}

	// 处理收到的消息
	// 从队列中接收消息并处理
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者名称
		true,   // 自动回复确认
		false,  // 独占消费者
		false,  // 没有自动确认
		false,  // 额外参数
		nil,    // 额外参数
	)

	if err != nil {
		return errno.NewErrNo("注册消费者失败！" + err.Error())
	}
	for d := range msgs {
		// 将消息反序列化为 Comment 结构体
		var comment model.Comment
		if err := json.Unmarshal(d.Body, &comment); err != nil {
			return errno.NewErrNo("解析JSON数据失败！" + err.Error())
		}
		// 保存评论到数据库
		if err := db.CreateComment(ctx, &comment); err != nil {
			return errno.NewErrNo("消息队列保存到数据库失败！" + err.Error())
		}
		redis.IncrVideoField(ctx, comment.VideoId, "comment_count", 1)
	}
	return nil
}
