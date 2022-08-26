package main

import (
	"context"

	"conversation-bot/config"
	"conversation-bot/internal/service/controllerservice"
	"conversation-bot/pkg/container"
	"conversation-bot/pkg/gpooling"
	handleossignal "conversation-bot/pkg/handle-os-signal"
	"conversation-bot/pkg/l"
	"conversation-bot/pkg/telegram"
	validator "conversation-bot/pkg/validator"
	"github.com/zelenin/go-tdlib/client"
)

func bootstrap(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	_, cancel := context.WithCancel(context.Background())
	shutdown.HandleDefer(cancel)

	container.NamedSingleton("gpooling", func() gpooling.IPool {
		return gpooling.New(cfg.MaxPoolSize, ll)
	})

	container.NamedSingleton("validator", func() validator.IValidator {
		return validator.New()
	})

	//region init store
	//endregion

	//region init agent
	teleConfigs := cfg.BotConfig.TelegramConfigs
	tdClients := make([]*client.Client, 0, len(teleConfigs))
	for _, tlConfig := range teleConfigs {
		tdClient := telegram.New(tlConfig)
		tdClients = append(tdClients, tdClient)
	}
	container.NamedSingleton("tdClients", func() []*client.Client {
		return tdClients
	})
	container.NamedSingleton("listener", func() *client.Client {
		return tdClients[0]
	})

	for i, conversationChannel := range cfg.BotConfig.ConversationChannels {
		if conversationChannel.ID == 0 {
			chat, err := tdClients[0].SearchChatsOnServer(&client.SearchChatsOnServerRequest{
				Query: conversationChannel.Name,
				Limit: 1,
			})
			if err != nil || len(chat.ChatIds) == 0 {
				ll.Error("cannot find chat with name", l.String("chat_name", conversationChannel.Name), l.Error(err))
			} else {
				ll.Info("found channel with name", l.String("chat_name", conversationChannel.Name), l.Int64("chat_id", chat.ChatIds[0]))
				cfg.BotConfig.ConversationChannels[i].ID = chat.ChatIds[0]
			}
		}
	}
	//endregion

	//region init service
	container.NamedSingleton("controllerService", func() controllerservice.IService {
		return controllerservice.New(cfg.BotConfig)
	})
	//endregion
}
