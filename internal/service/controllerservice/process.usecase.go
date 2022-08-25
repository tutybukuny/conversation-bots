package controllerservice

import (
	"context"
	"errors"
	"math/rand"

	"github.com/zelenin/go-tdlib/client"

	"conversation-bot/pkg/l"
)

func (s *serviceImpl) Process(ctx context.Context, message *client.Message) error {
	if message.ChatId != s.controlChannelID {
		s.ll.Debug("ignore which is not control channel")
		return nil
	}

	s.ll.Info("received message from control channel", l.Object("message", message))

	idx := 0
	for {
		idx = rand.Intn(len(s.tdClients))
		if idx == s.lastSenderIdx {
			continue
		}
		s.lastSenderIdx = idx
		break
	}

	sender := s.tdClients[idx]
	msg := &client.SendMessageRequest{}

	switch message.Content.MessageContentType() {
	case client.TypeMessageAnimation:
		content := message.Content.(*client.MessageAnimation)
		msg.InputMessageContent = &client.InputMessageAnimation{
			Animation: &client.InputFileRemote{Id: content.Animation.Animation.Remote.Id},
			Caption:   content.Caption,
		}
	case client.TypeMessageAudio:
		content := message.Content.(*client.MessageAudio)
		msg.InputMessageContent = &client.InputMessageAudio{
			Audio:   &client.InputFileRemote{Id: content.Audio.Audio.Remote.Id},
			Caption: content.Caption,
		}
	case client.TypeMessagePhoto:
		content := message.Content.(*client.MessagePhoto)
		msg.InputMessageContent = &client.InputMessagePhoto{
			Photo:   &client.InputFileRemote{Id: content.Photo.Sizes[len(content.Photo.Sizes)-1].Photo.Remote.Id},
			Caption: content.Caption,
		}
	case client.TypeMessageVideo:
		content := message.Content.(*client.MessageVideo)
		msg.InputMessageContent = &client.InputMessageVideo{
			Video:   &client.InputFileRemote{Id: content.Video.Video.Remote.Id},
			Caption: content.Caption,
		}
	case client.TypeMessageText:
		content := message.Content.(*client.MessageText)
		msg.InputMessageContent = &client.InputMessageText{Text: content.Text}
	default:
		return errors.New("not handled message type")
	}

	for _, conversationChannelID := range s.conversationChannelIDs {
		msg.ChatId = conversationChannelID
		sentMessage, err := sender.SendMessage(msg)
		if err != nil {
			s.ll.Error("error when sending message", l.Int64("conversation_channel_id", conversationChannelID), l.Object("msg", msg), l.Error(err))
			return err
		}
		s.ll.Info("sent message", l.Object("msg", msg), l.Object("sent_message", sentMessage))
	}

	return nil
}
