package kinopoisk

type Response struct {
	Docs []Document
}

type Genre struct {
	Name string
}

type Country struct {
	Name string
}

type PosterURL struct {
	URL string
}

type Rating struct {
	KP float32
}

type Document struct {
	ID          int
	Name        string
	Year        int
	Description string `json:"shortDescription"`
	Length      int    `json:"movieLength"`
	Poster      PosterURL
	Rating      Rating
}
