package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"log"
	"time"

	"github.com/crossle/imeos-news-mixin-bot/config"
	"github.com/crossle/imeos-news-mixin-bot/durable"
	"github.com/crossle/imeos-news-mixin-bot/models"
	"github.com/crossle/imeos-news-mixin-bot/session"

	bot "github.com/MixinNetwork/bot-api-go-client"
)

type BlazeService struct {
}

type ResponseMessage struct {
	blazeClient *bot.BlazeClient
}

func (service *BlazeService) Run(ctx context.Context) error {
	for {
		blazeClient := bot.NewBlazeClient(config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
		response := ResponseMessage{
			blazeClient: blazeClient,
		}
		if err := blazeClient.Loop(ctx, response); err != nil {
			session.Logger(ctx).Error(err)
		}
		session.Logger(ctx).Info("connection loop end")
		time.Sleep(time.Second)
	}
	return nil
}

func StartBlaze(db *sql.DB) error {
	log.Println("start blaze")
	logger, err := durable.NewLoggerClient("", true)
	if err != nil {
		return err
	}
	defer logger.Close()
	ctx, cancel := newBlazeContext(db, logger)
	defer cancel()

	b := BlazeService{}
	b.Run(ctx)
	return nil
}

func (r ResponseMessage) OnMessage(ctx context.Context, msg bot.MessageView, uid string) error {
	if msg.Category != bot.MessageCategorySystemAccountSnapshot && msg.Category != bot.MessageCategorySystemConversation && msg.ConversationId == bot.UniqueConversationId(config.MixinClientId, msg.UserId) {
		if msg.Category == "PLAIN_TEXT" {
			data, err := base64.StdEncoding.DecodeString(msg.Data)
			if err != nil {
				return bot.BlazeServerError(ctx, err)
			}
			if "/start" == string(data) {
				_, err = models.CreateSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := r.blazeClient.SendPlainText(ctx, msg, "订阅成功"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else if "/stop" == string(data) {
				err = models.RemoveSubscriber(ctx, msg.UserId)
				if err == nil {
					if err := r.blazeClient.SendPlainText(ctx, msg, "已取消订阅"); err != nil {
						return bot.BlazeServerError(ctx, err)
					}
				}
			} else {
				content := `发送 /start 订阅消息
发送 /stop 取消订阅`
				if err := r.blazeClient.SendPlainText(ctx, msg, content); err != nil {
					return bot.BlazeServerError(ctx, err)
				}
			}
		} else {
			content := `发送 /start 订阅消息
发送 /stop 取消订阅`
			if err := r.blazeClient.SendPlainText(ctx, msg, content); err != nil {
				return bot.BlazeServerError(ctx, err)
			}
		}
	}
	return nil
}

func newBlazeContext(db *sql.DB, client *durable.LoggerClient) (context.Context, context.CancelFunc) {
	ctx := session.WithLogger(context.Background(), durable.BuildLogger(client, "blaze", nil))
	ctx = session.WithDatabase(ctx, db)
	return context.WithCancel(ctx)
}
