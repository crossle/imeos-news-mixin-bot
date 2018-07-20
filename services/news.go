package services

import (
	"context"
	"encoding/base64"
	"log"

	bot "github.com/MixinNetwork/bot-api-go-client"
	"github.com/crossle/imeos-news-mixin-bot/config"
	"github.com/crossle/imeos-news-mixin-bot/models"
	"github.com/jasonlvhit/gocron"
)

type NewsService struct{}

type Stats struct {
	prevStoryId int64
}

func (self *Stats) getPrevTopStoryId() int64 {
	return self.prevStoryId
}

func (self *Stats) updatePrevTopStoryId(id int64) {
	self.prevStoryId = id
}

func getTopStory() NewsFlash {
	stories, err := GetStories()
	if err != nil {
		return NewsFlash{}
	}
	return stories[len(stories)-1]
}

func sendTopStoryToChannel(ctx context.Context, stats *Stats) {
	prevStoryId := stats.getPrevTopStoryId()
	stories, _ := GetStories()
	for i := len(stories) - 1; i >= 0; i-- {
		story := stories[i]
		if story.IssueTime > prevStoryId {
			log.Printf("Sending top story to channel...")
			stats.updatePrevTopStoryId(story.IssueTime)
			subscribers, _ := models.FindSubscribers(ctx)
			for _, subscriber := range subscribers {
				conversationId := bot.UniqueConversationId(config.MixinClientId, subscriber.UserId)
				data := base64.StdEncoding.EncodeToString([]byte(story.Content))
				bot.PostMessage(ctx, conversationId, subscriber.UserId, bot.NewV4().String(), "PLAIN_TEXT", data, config.MixinClientId, config.MixinSessionId, config.MixinPrivateKey)
			}
		} else {
			log.Printf("Same top story ID: %d, no message sent.", prevStoryId)
		}
	}
}
func (service *NewsService) Run(ctx context.Context) error {
	topStory := getTopStory()
	stats := &Stats{topStory.IssueTime}
	gocron.Every(5).Minutes().Do(sendTopStoryToChannel, ctx, stats)
	<-gocron.Start()
	return nil
}
