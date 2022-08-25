package controllerservice

import (
	"context"

	"github.com/zelenin/go-tdlib/client"

	"conversation-bot/config"
	"conversation-bot/pkg/container"
	"conversation-bot/pkg/gpooling"
	"conversation-bot/pkg/l"
)

type IService interface {
	Process(ctx context.Context, message *client.Message) error
}

type serviceImpl struct {
	ll        l.Logger         `container:"name"`
	gpooling  gpooling.IPool   `container:"name"`
	tdClients []*client.Client `container:"name"`

	controlChannelID       int64
	conversationChannelIDs []int64
	lastSenderIdx          int
}

func New(cfg *config.BotConfig) *serviceImpl {
	service := &serviceImpl{
		controlChannelID:       cfg.ControlChannelID,
		conversationChannelIDs: make([]int64, 0, len(cfg.ConversationChannels)),
	}
	container.Fill(service)

	for _, channel := range cfg.ConversationChannels {
		service.conversationChannelIDs = append(service.conversationChannelIDs, channel.ID)
	}

	return service
}
