package telegram

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"kinopoisk-telegram-bot/pkg/clients/kinopoisk"
	"kinopoisk-telegram-bot/pkg/storage"
	"strconv"
	"strings"
)

const (
	listPrefix   = "list_"
	deletePrefix = "delete_"
)

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Start)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Help)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleListCommand(message *tgbotapi.Message) error {
	movies, err := b.storage.GetAll(message.Chat.ID, storage.FavoriteMovies)
	if err != nil {
		return err
	}
	if movies == nil {
		return errFavListIsEmpty
	}
	err = b.sendFavMovies(message.Chat.ID, movies, listPrefix)
	if err != nil {
		return err
	}
	b.userState.SetAwaitingResponse(message.Chat.ID, true)
	return nil
}

func (b *Bot) handleDeleteCommand(message *tgbotapi.Message) error {
	movies, err := b.storage.GetAll(message.Chat.ID, storage.FavoriteMovies)
	if err != nil {
		return err
	}
	if movies == nil {
		return errFavListIsEmpty
	}
	err = b.sendFavMovies(message.Chat.ID, movies, deletePrefix)
	if err != nil {
		return err
	}
	b.userState.SetAwaitingResponse(message.Chat.ID, true)
	return nil
}

func (b *Bot) saveTempMovies(message *tgbotapi.Message) error {
	movies, err := b.client.Get(message.Text)
	if err != nil {
		return errCantFetchMovies
	}
	for _, movie := range movies.Docs {
		if err = b.storage.Save(message.Chat.ID, movie, storage.TemporaryMovies); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) getAndSendInfo(chatID int64) error {
	movie, err := b.storage.Get(chatID, storage.TemporaryMovies)
	if err != nil {
		b.userState.SetAwaitingResponse(chatID, false)
		if errors.Is(err, sql.ErrNoRows) {
			return errEndOfSearch
		}
		return err
	}
	if err = b.storage.Remove(chatID, movie.ID, storage.TemporaryMovies); err != nil {
		return err
	}
	if err = b.sendInfo(chatID, movie); err != nil {
		return err
	}
	err = b.sendButtons(chatID, b.messages.IsTheMovie)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) sendInfo(chatID int64, movie *kinopoisk.Document) error {
	if movie.Poster.URL != "" {
		if err := b.withPoster(chatID, movie); err != nil {
			return err
		}
	} else {
		if err := b.withoutPoster(chatID, movie); err != nil {
			return err
		}
	}
	b.userState.SetFavoriteMovie(chatID, *movie)
	return nil
}

func (b *Bot) withPoster(chatID int64, movie *kinopoisk.Document) error {
	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(movie.Poster.URL))
	msg.Caption = fmt.Sprintf(
		"\"%s\"\n\nГод: %d\nОписание: %s\nДлительность: %d мин.\nРейтинг кинопоиска: %.2f\n",
		movie.Name, movie.Year, movie.Description,
		movie.Length, movie.Rating.KP,
	)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) withoutPoster(chatID int64, movie *kinopoisk.Document) error {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(
		"*Постер отсутствует*\n\n\"%s\"\n\nГод: %d\nОписание: %s\nДлительность: %d мин.\nРейтинг кинопоиска: %.2f\n",
		movie.Name, movie.Year, movie.Description,
		movie.Length, movie.Rating.KP))
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) sendButtons(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	var yes, no string
	if b.userState.GetButtonState(chatID) {
		yes = "yes_add"
		no = "no_add"
	} else {
		yes = "yes_confirm"
		no = "no_confirm"
	}
	yesButton := tgbotapi.NewInlineKeyboardButtonData("Да", yes)
	noButton := tgbotapi.NewInlineKeyboardButtonData("Нет", no)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yesButton, noButton))
	msg.ReplyMarkup = keyboard

	_, err := b.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("cant send buttons: %w", err)
	}
	b.userState.SetAwaitingResponse(chatID, true)
	return nil
}

func (b *Bot) handleMovieButtonListCase(callback *tgbotapi.CallbackQuery) error {
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	b.userState.SetAwaitingResponse(callback.Message.Chat.ID, false)
	movieID, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "list_"))
	movie, err := b.storage.GetByID(callback.Message.Chat.ID, movieID, storage.FavoriteMovies)
	if err != nil {
		return err
	}
	if err = b.sendInfo(callback.Message.Chat.ID, movie); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleMovieButtonDeleteCase(callback *tgbotapi.CallbackQuery) error {
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	b.userState.SetAwaitingResponse(callback.Message.Chat.ID, false)
	movieID, _ := strconv.Atoi(strings.TrimPrefix(callback.Data, "delete_"))
	if err := b.storage.Remove(callback.Message.Chat.ID, movieID, storage.FavoriteMovies); err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, b.messages.DeleteSuccessfully)
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleLeft(message *tgbotapi.Message) error {
	pageNum := b.userState.GetPageNum(message.Chat.ID)
	if pageNum == 0 {
		return nil
	}
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	b.userState.SetPageNum(message.Chat.ID, pageNum-1)
	if err := b.handleListCommand(message); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleRight(message *tgbotapi.Message) error {
	pageNum := b.userState.GetPageNum(message.Chat.ID)
	movies, err := b.storage.GetAll(message.Chat.ID, storage.FavoriteMovies)
	if err != nil {
		return err
	}
	if pageNum > len(movies)/moviesOnPage-1 {
		return nil
	}
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	b.userState.SetPageNum(message.Chat.ID, pageNum+1)
	err = b.sendFavMovies(message.Chat.ID, movies, listPrefix)
	if err != nil {
		return err
	}
	b.userState.SetAwaitingResponse(message.Chat.ID, true)
	return nil
}

func (b *Bot) handleYesAdd(message *tgbotapi.Message) error {
	b.userState.SetAwaitingResponse(message.Chat.ID, false)
	b.userState.SetButtonState(message.Chat.ID, false)
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	if err := b.storage.RemoveAll(message.Chat.ID, storage.TemporaryMovies); err != nil {
		return err
	}
	if err := b.storage.Save(message.Chat.ID, b.userState.GetFavoriteMovie(message.Chat.ID), storage.FavoriteMovies); err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AddSuccessfully)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleNoAdd(message *tgbotapi.Message) error {
	b.userState.SetAwaitingResponse(message.Chat.ID, false)
	b.userState.SetButtonState(message.Chat.ID, false)
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	if err := b.storage.RemoveAll(message.Chat.ID, storage.TemporaryMovies); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleNoConfirm(message *tgbotapi.Message) error {
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	err := b.getAndSendInfo(message.Chat.ID)
	if err != nil {
		return fmt.Errorf("cant send a message: %w", err)
	}
	return nil
}

func (b *Bot) handleYesConfirm(message *tgbotapi.Message) error {
	b.userState.SetButtonState(message.Chat.ID, true)
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)); err != nil || !resp.Ok {
		return err
	}
	if err := b.sendButtons(message.Chat.ID, b.messages.WantToSave); err != nil {
		return err
	}
	return nil
}

func (b *Bot) sendFavMovies(chatID int64, movies []kinopoisk.Document, prefix string) error {
	rows := b.createRows(chatID, movies, prefix)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, b.messages.YourMovies)
	msg.ReplyMarkup = keyboard
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) createRows(chatID int64, movies []kinopoisk.Document, prefix string) [][]tgbotapi.InlineKeyboardButton {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton
	currentPage := b.userState.GetPageNum(chatID)
	numPages := len(movies)/moviesOnPage + 1
	pageBegin := moviesOnPage * currentPage
	var pageEnd int
	if len(movies)-1 < pageBegin+moviesOnPage {
		pageEnd = len(movies)
	} else {
		pageEnd = pageBegin + moviesOnPage
	}
	for i, movie := range movies[pageBegin:pageEnd] {
		if movie.Name != "" {
			emoji := getEmoji(i)
			newButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%v \"%v\" %d", emoji,
				movie.Name, movie.Year), prefix+strconv.Itoa(movie.ID))
			row = append(row, newButton)
			if len(row) == 2 || (i == len(movies)%moviesOnPage-1 && currentPage == numPages-1) {
				rows = append(rows, row)
				row = nil
			}
		}
	}

	leftButton := tgbotapi.NewInlineKeyboardButtonData("<", "left")
	pageNumButton := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", currentPage+1, numPages), "pagenum")
	rightButton := tgbotapi.NewInlineKeyboardButtonData(">", "right")
	navigationRow := []tgbotapi.InlineKeyboardButton{leftButton, pageNumButton, rightButton}
	rows = append(rows, navigationRow)
	return rows
}

func getEmoji(index int) string {
	emojis := []string{"🎥", "🎬", "🍿", "🎞", "📽"}
	return emojis[index%5]
}
