package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"kinopoisk-telegram-bot/pkg/clients/kinopoisk"
	"kinopoisk-telegram-bot/pkg/clients/telegram"
	"kinopoisk-telegram-bot/pkg/config"
	"kinopoisk-telegram-bot/pkg/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func Init(cfg *config.Config) (db *sql.DB, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error creating db: %w", err)
		}
	}()
	db, err = sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS ` + string(storage.TemporaryMovies) + `;
	CREATE TABLE IF NOT EXISTS ` + string(storage.TemporaryMovies) + ` (
	user_id INTEGER,
    id INTEGER,
    movie_name TEXT,
    movie_year INTEGER,
    movie_description TEXT,
    movie_length INTEGER,
    movie_poster TEXT,
    movie_rating REAL)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + string(storage.FavoriteMovies) + ` (
	user_id INTEGER,
    id INTEGER,
    movie_name TEXT,
    movie_year INTEGER,
    movie_description TEXT,
    movie_length INTEGER,
    movie_poster TEXT,
    movie_rating REAL)`)
	if err != nil {
		return nil, err
	}
	return db, err
}

func (s *Storage) Save(userID int64, movie kinopoisk.Document, table storage.Table) error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM "+string(table)+" WHERE user_id=? AND id=?", userID, movie.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return telegram.ErrRecExists
	}

	_, err = s.db.Exec("INSERT INTO "+string(table)+" (user_id, id, movie_name, movie_year, movie_description, movie_length, movie_poster, movie_rating) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		userID, movie.ID, movie.Name,
		movie.Year, movie.Description,
		movie.Length, movie.Poster.URL,
		movie.Rating.KP,
	)
	if err != nil {
		return telegram.ErrCantInsert
	}
	return nil
}

func (s *Storage) Get(userID int64, table storage.Table) (*kinopoisk.Document, error) {
	row := s.db.QueryRow("SELECT id, movie_name, movie_year, movie_description, movie_length, movie_poster, movie_rating FROM "+string(table)+" WHERE user_id=? LIMIT 1",
		userID)

	var movie kinopoisk.Document
	err := row.Scan(&movie.ID, &movie.Name, &movie.Year, &movie.Description, &movie.Length, &movie.Poster.URL, &movie.Rating.KP)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *Storage) GetAll(userID int64, table storage.Table) ([]kinopoisk.Document, error) {
	rows, err := s.db.Query("SELECT id, movie_name, movie_year, movie_description, movie_length, movie_poster, movie_rating FROM "+string(table)+" WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []kinopoisk.Document
	var movie kinopoisk.Document
	for rows.Next() {
		err := rows.Scan(&movie.ID, &movie.Name, &movie.Year, &movie.Description, &movie.Length, &movie.Poster.URL, &movie.Rating.KP)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (s *Storage) GetByID(userID int64, movieID int, table storage.Table) (*kinopoisk.Document, error) {
	row := s.db.QueryRow("SELECT id, movie_name, movie_year, movie_description, movie_length, movie_poster, movie_rating FROM "+string(table)+" WHERE id=? AND user_id=? LIMIT 1",
		movieID, userID)

	var movie kinopoisk.Document
	err := row.Scan(&movie.ID, &movie.Name, &movie.Year, &movie.Description, &movie.Length, &movie.Poster.URL, &movie.Rating.KP)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *Storage) Remove(userID int64, movieID int, table storage.Table) error {
	_, err := s.db.Exec("DELETE FROM "+string(table)+" WHERE user_id=? AND id=?", userID, movieID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) RemoveAll(userID int64, table storage.Table) error {
	_, err := s.db.Exec("DELETE FROM "+string(table)+" WHERE user_id=?", userID)
	if err != nil {
		return err
	}
	return nil
}
