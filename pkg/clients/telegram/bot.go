package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kinopoisk-telegram-bot/pkg/clients/kinopoisk"
	"kinopoisk-telegram-bot/pkg/config"
	"kinopoisk-telegram-bot/pkg/storage"
	"kinopoisk-telegram-bot/pkg/userstate"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	client    *kinopoisk.Client
	storage   storage.Storage
	userState *userstate.UserState
	messages  config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, client *kinopoisk.Client, storage storage.Storage, userState *userstate.UserState, messages config.Messages) *Bot {
	return &Bot{
		bot:       bot,
		client:    client,
		storage:   storage,
		userState: userState,
		messages:  messages,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				if err := b.handleCommand(update.Message); err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
				continue
			}

			if err := b.handleMessage(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
		} else if update.CallbackQuery != nil {
			err := b.handleCallback(update.CallbackQuery)
			if err != nil {
				b.handleError(update.CallbackQuery.Message.Chat.ID, err)
			}
		}
	}
	return nil
}
