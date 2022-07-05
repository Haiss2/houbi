package bot

import (
	"log"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	botAPI "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	telegramTokenFlag        = "telegram-token"
	telegramChatIDFlag       = "telegram-chat-id"
	MarkdownMode             = "markdown"
	TELEGRAM_BOT_TOKEN       = "5001043957:AAF1uqnfHzq8uRa1Gejswd0UHzIIPkMjilc"
	TELEGRAM_CHAT_ID   int64 = -732240513
)

type TelegramBot struct {
	api    *botAPI.BotAPI
	chatId int64
}

func NewTelegramBot() (*TelegramBot, error) {
	api, err := botAPI.NewBotAPI(TELEGRAM_BOT_TOKEN)

	if err != nil {
		return nil, err
	}
	chatId := TELEGRAM_CHAT_ID
	err = validation.Validate(chatId, validation.Required)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{
		api:    api,
		chatId: chatId,
	}, nil
}

func (t *TelegramBot) Notify(msg, fmt string) error {
	return t.NotifyToGroup(t.chatId, msg, fmt)
}

func (t *TelegramBot) NotifyToGroup(chatId int64, msg, fmt string) error {
	sendMsg := botAPI.NewMessage(chatId, msg)
	sendMsg.ParseMode = parseMode(fmt)
	_, err := t.api.Send(sendMsg)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (t *TelegramBot) NotifyWithTag(msg, tags, fmt string) error {
	return t.NotifyToGroupWithTags(t.chatId, msg, tags, fmt)
}

func (t *TelegramBot) NotifyToGroupWithTags(chatId int64, msg, tags, fmt string) error {
	sendMsg := botAPI.NewMessage(chatId, msg+"\n"+tags)
	sendMsg.ParseMode = parseMode(fmt)
	_, err := t.api.Send(sendMsg)
	if err != nil {
		log.Print(err)
	}
	return err
}

func parseMode(fmt string) string {
	switch fmt {
	case MarkdownMode:
		return MarkdownMode
	default:
		return ""
	}
}
