package userstate

import (
	"kinopoisk-telegram-bot/pkg/clients/kinopoisk"
	"sync"
)

type UserState struct {
	mu               sync.Mutex
	awaitingResponse map[int64]bool
	favoriteMovie    map[int64]kinopoisk.Document
	pageNum          map[int64]int
}

func NewUserState() *UserState {
	return &UserState{
		awaitingResponse: make(map[int64]bool),
		favoriteMovie:    make(map[int64]kinopoisk.Document),
		pageNum:          make(map[int64]int),
	}
}

func (us *UserState) SetAwaitingResponse(chatID int64, awaiting bool) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.awaitingResponse[chatID] = awaiting
}

func (us *UserState) IsAwaitingResponse(chatID int64) bool {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.awaitingResponse[chatID]
}

func (us *UserState) SetFavoriteMovie(chatID int64, movie kinopoisk.Document) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.favoriteMovie[chatID] = movie
}

func (us *UserState) GetFavoriteMovie(chatID int64) kinopoisk.Document {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.favoriteMovie[chatID]
}

func (us *UserState) SetPageNum(chatID int64, index int) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.pageNum[chatID] = index
}

func (us *UserState) GetPageNum(chatID int64) int {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.pageNum[chatID]
}
