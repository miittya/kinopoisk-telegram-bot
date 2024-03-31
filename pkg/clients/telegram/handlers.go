package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCmd     = "start"
	helpCmd      = "help"
	listCmd      = "list"
	moviesOnPage = 6
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case startCmd:
		return b.handleStartCommand(message)
	case helpCmd:
		return b.handleHelpCommand(message)
	case listCmd:
		return b.handleListCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	if b.userState.IsAwaitingResponse(message.Chat.ID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.PushTheButton)
		if _, err := b.bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
	if err := b.saveTempMovies(message); err != nil {
		return err
	}
	if err := b.getAndSendInfo(message.Chat.ID); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	switch callback.Data {
	case "yes_add":
		if err := b.handleYesAdd(callback.Message); err != nil {
			return err
		}
	case "no_add":
		if err := b.handleNoAdd(callback.Message); err != nil {
			return err
		}
	case "yes_confirm":
		if err := b.handleYesConfirm(callback.Message); err != nil {
			return err
		}
	case "no_confirm":
		if err := b.handleNoConfirm(callback.Message); err != nil {
			return err
		}
	case "left":
		if err := b.handleLeft(callback.Message); err != nil {
			return err
		}
	case "right":
		if err := b.handleRight(callback.Message); err != nil {
			return err
		}
	case "pagenum":
		return nil
	default:
		if err := b.handleMovieButton(callback); err != nil {
			return err
		}
	}
	return nil
}
