package dep

import (
	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/scrapper/mangasail"
	"github.com/bigscreen/manga-scrapper/service"
)

type AppDependencies struct {
	FetchServiceMap map[common.FetchServiceKey]service.FetchService
}

func InitAppDependencies() AppDependencies {
	fetchSvcMap := map[common.FetchServiceKey]service.FetchService{
		common.FSKeyMangasail: service.NewMangsailService(service.MangasailServiceParams{
			HomeScrapper:           mangasail.NewHomePageScrapper(),
			MangaDetailsScrapper:   mangasail.NewMangaDetailsPageScrapper(),
			ChapterDetailsScrapper: mangasail.NewChapterDetailsPageScrapper(),
		}),
	}

	return AppDependencies{
		FetchServiceMap: fetchSvcMap,
	}
}
