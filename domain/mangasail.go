package domain

import "encoding/json"

type (
	MangasailMangas   []MangasailManga
	MangasailChapters []MangasailChapter
	MangasailImages   []MangasailImage
)

type MangasailManga struct {
	Name        string            `json:"name,omitempty"`
	Path        string            `json:"path,omitempty"`
	IconURL     string            `json:"icon_url,omitempty"`
	Author      string            `json:"author,omitempty"`
	Status      string            `json:"status,omitempty"`
	ReleaseYear string            `json:"release_year,omitempty"`
	Description string            `json:"description,omitempty"`
	Genres      []string          `json:"genres,omitempty"`
	Chapters    MangasailChapters `json:"chapters,omitempty"`
}

type MangasailChapter struct {
	Title        string          `json:"title,omitempty"`
	Path         string          `json:"path,omitempty"`
	LastModified string          `json:"last_modified,omitempty"`
	Images       MangasailImages `json:"images,omitempty"`
}

type MangasailImage struct {
	ID       string `json:"id,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type MangasailHomeMangas struct {
	DailyHotMangas MangasailMangas `json:"daily_hot_mangas,omitempty"`
	LatestMangas   MangasailMangas `json:"latest_mangas,omitempty"`
	HotMangas      MangasailMangas `json:"hot_mangas,omitempty"`
	NewMangas      MangasailMangas `json:"new_mangas,omitempty"`
}

func (m MangasailManga) ToString() string {
	return getJson(m)
}

func (m MangasailChapter) ToString() string {
	return getJson(m)
}

func (m MangasailImage) ToString() string {
	return getJson(m)
}

func (m MangasailHomeMangas) ToString() string {
	return getJson(m)
}

func getJson(i interface{}) string {
	j, err := json.Marshal(i)
	if err == nil {
		return string(j)
	}
	return ""
}
