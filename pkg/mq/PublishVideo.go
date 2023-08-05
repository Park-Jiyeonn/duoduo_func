package mq

import (
	"context"
	"duoduo_fun/dal/db"
	"duoduo_fun/dal/db/model"
	"duoduo_fun/pkg/consts"
	util "duoduo_fun/util/ffmpeg"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func Produce(uid int, playUrl, coverUrl string) error {
	conn, err := amqp.Dial(consts.RabbitMQDSN)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	msg := &model.Video{
		UserId:   uid,
		CoverUrl: coverUrl,
		PlayUrl:  playUrl,
	}
	videoJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(
		consts.VideoExchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",
		consts.VideoExchange,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        videoJSON,
		})
	if err != nil {
		return err
	}
	return nil
}

func Consume(ctx context.Context) {
	conn, err := amqp.Dial(consts.RabbitMQDSN)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		consts.VideoExchange,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	err = ch.QueueBind(
		queue.Name,
		"upload",
		consts.VideoExchange,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}
	for msg := range msgs {
		err = msg.Ack(true)
		if err != nil {
			log.Println(err)
		}
		go func(msg amqp.Delivery) {
			videoJSON := msg.Body
			var publishInfo model.Video
			err := json.Unmarshal(videoJSON, &publishInfo)
			if err != nil {
				log.Println(err)
			}
			err = util.Cover(publishInfo.PlayUrl, publishInfo.CoverUrl)
			if err != nil {
				//err = Produce(publishInfo.UserId, publishInfo.Title, publishInfo.PlayUrl)
				if err != nil {
					log.Println(err)
				}
				return
			}
			err = db.CreateVideo(ctx, publishInfo.PlayUrl, publishInfo.CoverUrl, publishInfo.Title, publishInfo.UserId)
			if err != nil {
				//err = Produce(publishInfo.UserId, publishInfo.Title, publishInfo.PlayUrl)
				if err != nil {
					log.Println(err)
				}
				return
			}
		}(msg)
	}
}
