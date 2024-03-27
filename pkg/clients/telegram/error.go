package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errFavListIsEmpty  = errors.New("list is empty")
	errCantFetchMovies = errors.New("error fetching movies")
	ErrRecExists       = errors.New("record already exists")
	ErrCantInsert      = errors.New("cant insert movie")
	errEndOfSearch     = errors.New("end of search")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)
	switch {
	case errors.Is(err, errFavListIsEmpty):
		msg.Text = b.messages.FavListIsEmpty
	case errors.Is(err, errCantFetchMovies):
		msg.Text = b.messages.CantFetchMovies
	case errors.Is(err, ErrRecExists):
		msg.Text = b.messages.RecExists
	case errors.Is(err, ErrCantInsert):
		msg.Text = b.messages.CantInsert
	case errors.Is(err, errEndOfSearch):
		msg.Text = b.messages.EndOfSearch
	}
	b.bot.Send(msg)
}
