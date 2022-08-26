package controllerservice

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/zelenin/go-tdlib/client"

	"conversation-bot/pkg/l"
)

func (s *serviceImpl) Process(ctx context.Context, message *client.Message) error {
	if message.ChatId != s.controlChannelID {
		s.ll.Debug("ignore which is not control channel")
		return nil
	}

	s.ll.Info("received message from control channel", l.Object("message", message))

	oriCaption, err := s.getCaption(message)
	if err != nil {
		return err
	}

	idx, caption := s.getBotIdx(oriCaption.Text)
	if idx < 0 || idx >= len(s.tdClients) {
		for {
			idx = rand.Intn(len(s.tdClients))
			if idx == s.lastSenderIdx {
				continue
			}
			break
		}
	}
	s.lastSenderIdx = idx

	sender := s.tdClients[idx]
	msg := &client.SendMessageRequest{}

	switch message.Content.MessageContentType() {
	case client.TypeMessageAnimation:
		content := message.Content.(*client.MessageAnimation)
		msg.InputMessageContent = &client.InputMessageAnimation{
			Animation: &client.InputFileRemote{Id: content.Animation.Animation.Remote.Id},
			Caption:   &client.FormattedText{Text: caption},
		}
	case client.TypeMessageAudio:
		content := message.Content.(*client.MessageAudio)
		msg.InputMessageContent = &client.InputMessageAudio{
			Audio:   &client.InputFileRemote{Id: content.Audio.Audio.Remote.Id},
			Caption: &client.FormattedText{Text: caption},
		}
	case client.TypeMessagePhoto:
		content := message.Content.(*client.MessagePhoto)
		msg.InputMessageContent = &client.InputMessagePhoto{
			Photo:   &client.InputFileRemote{Id: content.Photo.Sizes[len(content.Photo.Sizes)-1].Photo.Remote.Id},
			Caption: &client.FormattedText{Text: caption},
		}
	case client.TypeMessageVideo:
		content := message.Content.(*client.MessageVideo)
		msg.InputMessageContent = &client.InputMessageVideo{
			Video:   &client.InputFileRemote{Id: content.Video.Video.Remote.Id},
			Caption: &client.FormattedText{Text: caption},
		}
	case client.TypeMessageText:
		msg.InputMessageContent = &client.InputMessageText{Text: &client.FormattedText{Text: caption}}
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

func (s *serviceImpl) getCaption(message *client.Message) (*client.FormattedText, error) {
	switch message.Content.MessageContentType() {
	case client.TypeMessageAnimation:
		content := message.Content.(*client.MessageAnimation)
		return content.Caption, nil
	case client.TypeMessageAudio:
		content := message.Content.(*client.MessageAudio)
		return content.Caption, nil
	case client.TypeMessagePhoto:
		content := message.Content.(*client.MessagePhoto)
		return content.Caption, nil
	case client.TypeMessageVideo:
		content := message.Content.(*client.MessageVideo)
		return content.Caption, nil
	case client.TypeMessageText:
		content := message.Content.(*client.MessageText)
		return content.Text, nil
	default:
		return nil, errors.New("not handled message type")
	}
}

func (s *serviceImpl) getBotIdx(text string) (int, string) {
	matches := s.botIdxRegex.FindStringSubmatch(text)
	if len(matches) == 0 {
		return -1, text
	}

	pos := s.botIdxRegex.FindStringSubmatchIndex(text)[0]
	idx, _ := strconv.Atoi(matches[1])
	return idx - 1, strings.TrimSpace(text[:pos] + text[pos+len(matches[0]):])
}
