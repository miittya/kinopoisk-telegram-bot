package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kinopoisk-telegram-bot/pkg/clients/kinopoisk"
	"kinopoisk-telegram-bot/pkg/clients/telegram"
	"kinopoisk-telegram-bot/pkg/config"
	"kinopoisk-telegram-bot/pkg/storage/sqlite"
	"kinopoisk-telegram-bot/pkg/userstate"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	userState := userstate.NewUserState()

	db, err := sqlite.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	storage := sqlite.NewStorage(db)

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	botAPI.Debug = true

	kpClient, err := kinopoisk.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	bot := telegram.NewBot(botAPI, kpClient, storage, userState, cfg.Messages)
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}
