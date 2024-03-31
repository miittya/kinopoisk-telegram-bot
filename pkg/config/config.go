package config

import (
	"github.com/caarlos0/env/v10"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TelegramToken       string `env:"TG_API_KEY"`
	KinopoiskToken      string `env:"X-API-KEY"`
	DBPath              string `yaml:"db_file"`
	APIHost             string
	EndPointMovieSearch string `yaml:"endpoint_movie_search"`
	Messages            Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start           string
	Help            string
	UnknownCommand  string `yaml:"unknown_command"`
	ListIsEmpty     string `yaml:"list_is_empty"`
	IsTheMovie      string `yaml:"is_the_movie"`
	WantToSave      string `yaml:"want_to_save"`
	AddSuccessfully string `yaml:"add_successfully"`
	YourMovies      string `yaml:"your_movies"`
	PushTheButton   string `yaml:"push_the_button"`
}

type Errors struct {
	Default         string `yaml:"default"`
	FavListIsEmpty  string `yaml:"fav_list_is_empty"`
	CantFetchMovies string `yaml:"cant_fetch_movies"`
	RecExists       string `yaml:"rec_exists"`
	CantInsert      string `yaml:"cant_insert"`
	EndOfSearch     string `yaml:"end_of_search"`
}

func Init() (*Config, error) {
	ymlFile, err := os.ReadFile("configs/main.yml")
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err = yaml.Unmarshal(ymlFile, cfg); err != nil {
		return nil, err
	}
	if err = env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
