package domain

import (
	"encoding/json"
)

type (
	Mangas   []Manga
	Chapters []Chapter
	Images   []Image
)

type Manga struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	IconURL     string   `json:"icon_url,omitempty"`
	Author      string   `json:"author,omitempty"`
	Status      string   `json:"status,omitempty"`
	ReleaseYear string   `json:"release_year,omitempty"`
	Description string   `json:"description,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Chapters    Chapters `json:"chapters,omitempty"`
}

type Chapter struct {
	ID           string `json:"id,omitempty"`
	Title        string `json:"title,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	Images       Images `json:"images,omitempty"`
}

type Image struct {
	ID       string `json:"id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type HomeMangas struct {
	DailyHotMangas Mangas `json:"daily_hot_mangas,omitempty"`
	LatestMangas   Mangas `json:"latest_mangas,omitempty"`
	PopularMangas  Mangas `json:"popular_mangas,omitempty"`
	NewMangas      Mangas `json:"new_mangas,omitempty"`
}

func (m Manga) ToString() string {
	return getJson(m)
}

func (m Chapter) ToString() string {
	return getJson(m)
}

func (m Image) ToString() string {
	return getJson(m)
}

func (m HomeMangas) ToString() string {
	return getJson(m)
}

func getJson(i interface{}) string {
	j, err := json.Marshal(i)
	if err == nil {
		return string(j)
	}
	return ""
}
