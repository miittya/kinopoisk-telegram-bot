package storage

import "kinopoisk-telegram-bot/pkg/clients/kinopoisk"

type Table string

const (
	TemporaryMovies Table = "temporary_movies"
	FavoriteMovies  Table = "favorite_movies"
)

type Storage interface {
	Save(userID int64, movie kinopoisk.Document, table Table) error
	Get(userID int64, table Table) (*kinopoisk.Document, error)
	GetAll(userID int64, table Table) ([]kinopoisk.Document, error)
	GetByID(userID int64, movieID int, table Table) (*kinopoisk.Document, error)
	Remove(userID int64, movieID int, table Table) error
	RemoveAll(userID int64, table Table) error
}
