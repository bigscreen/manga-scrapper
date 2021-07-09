package domain

type (
	Mangas   []Manga
	Chapters []Chapter
	Images   []Image
)

type Manga struct {
	ID          string
	Name        string
	IconURL     string
	Author      string
	Status      string
	ReleaseYear string
	Description string
	Genres      []string
	Chapters    Chapters
}

type Chapter struct {
	ID           string
	Title        string
	LastModified string
	Images       Images
}

type Image struct {
	ID       string
	ImageURL string
}

type HomeMangas struct {
	DailyHotMangas Mangas
	LatestMangas   Mangas
	PopularMangas  Mangas
	NewMangas      Mangas
}
