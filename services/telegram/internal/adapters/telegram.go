package adapters

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramAdapter(botToken string, chatID int64) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}
	return &TelegramAdapter{bot: bot, chatID: chatID}, nil
}

func (t *TelegramAdapter) SendMessage(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)

	_, err := t.bot.Send(msg)

	return err
}
