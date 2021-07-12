package common

type FetchServiceKey string

const (
	GetHomeCardsAPIPath      = "/mangajack/v1/home"
	GetMangaDetailsAPIPath   = "/mangajack/v1/manga"
	GetChapterDetailsAPIPath = "/mangajack/v1/chapter"

	ParamKeySource = "source"
	ParamKeyId     = "id"

	FSKeyMangasail FetchServiceKey = "mangasail"

	MangasailBaseURL    = "https://www.mangasail.co"
	MangasailPrefixPath = "/content/"
)
