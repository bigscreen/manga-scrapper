package contract

type RedirectType string
type HomeCardIdentifier string

const (
	RedirectToManga   RedirectType = "to_manga"
	RedirectToChapter RedirectType = "to_chapter"

	HomeCardDailyHotMangas HomeCardIdentifier = "daily_hot_mangas"
	HomeCardPopularMangas  HomeCardIdentifier = "popular_mangas"
	HomeCardLatestMangas   HomeCardIdentifier = "latest_mangas"
	HomeCardNewMangas      HomeCardIdentifier = "new_mangas"
)

type Home struct {
	PageTitle string     `json:"page_title"`
	Cards     []HomeCard `json:"cards"`
}

type HomeCard struct {
	Identifier HomeCardIdentifier `json:"identifier"`
	Label      string             `json:"label"`
	Content    []HomeCardContent  `json:"content"`
}

type HomeCardContent struct {
	Title       string              `json:"title"`
	SubTitle    string              `json:"sub_title,omitempty"`
	Info        string              `json:"info,omitempty"`
	IconURL     string              `json:"icon_url"`
	Redirection HomeCardRedirection `json:"redirection"`
}

type HomeCardRedirection struct {
	ID   string       `json:"id"`
	Type RedirectType `json:"type"`
}

func (m Home) String() string {
	return getJson(m)
}
