package listener

import (
	"context"

	"github.com/zelenin/go-tdlib/client"

	"conversation-bot/internal/service/controllerservice"
	"conversation-bot/pkg/container"
	"conversation-bot/pkg/gpooling"
	"conversation-bot/pkg/l"
)

type TelegramListener struct {
	ll                l.Logger                   `container:"name"`
	gpooling          gpooling.IPool             `container:"name"`
	listener          *client.Client             `container:"name"`
	controllerService controllerservice.IService `container:"name"`
}

func New() *TelegramListener {
	listener := &TelegramListener{}
	container.Fill(listener)

	return listener
}

func (tl *TelegramListener) Listen() {
	listener := tl.listener.GetListener()
	defer listener.Close()
	ctx := context.Background()

	for update := range listener.Updates {
		envelop, ok := update.(*client.UpdateNewMessage)
		if !ok {
			continue
		}
		message := envelop.Message

		if err := tl.controllerService.Process(ctx, message); err != nil {
			tl.ll.Error("cannot process message", l.Object("message", message), l.Error(err))
		}
	}
}
