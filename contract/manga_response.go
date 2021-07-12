package contract

type Manga struct {
	PageTitle      string         `json:"page_title"`
	HeaderImageURL string         `json:"header_image_url"`
	Info           MangaInfo      `json:"info"`
	Chapters       []MangaChapter `json:"chapters"`
}

type MangaInfo struct {
	Author      string   `json:"author,omitempty"`
	Status      string   `json:"status,omitempty"`
	ReleaseYear string   `json:"release_year,omitempty"`
	Description string   `json:"description,omitempty"`
	Genres      []string `json:"genres,omitempty"`
}

type MangaChapter struct {
	ID           string `json:"id,omitempty"`
	Title        string `json:"title,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
}
