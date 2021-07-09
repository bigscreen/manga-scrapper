package contract

type Chapter struct {
	PageTitle string   `json:"page_title"`
	ImageURLs []string `json:"image_urls"`
}

func (m Chapter) String() string {
	return getJson(m)
}
