package service

import (
	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/domain"
	"github.com/bigscreen/manga-scrapper/scrapper/mangasail"
)

type MangasailServiceParams struct {
	HomeScrapper           mangasail.HomePageScrapper
	MangaDetailsScrapper   mangasail.MangaDetailsPageScrapper
	ChapterDetailsScrapper mangasail.ChapterDetailsPageScrapper
}

type mangasailService struct {
	MangasailServiceParams
}

func NewMangsailService(params MangasailServiceParams) FetchService {
	return mangasailService{params}
}

func (m mangasailService) GetHomeCards() (contract.Home, error) {
	c, err := m.HomeScrapper.GetContent()
	if err != nil {
		return contract.Home{}, err
	}

	return contract.Home{
		PageTitle: "Mangajack Home",
		Cards: []contract.HomeCard{
			m.buildDailyHotMangasHomeCard(c.DailyHotMangas),
			m.buildPopularMangasHomeCard(c.PopularMangas),
			m.buildLatestMangasHomeCard(c.LatestMangas),
			m.buildNewMangasHomeCard(c.NewMangas),
		},
	}, nil
}

func (m mangasailService) GetMangaDetails(mangaId string) (contract.Manga, error) {
	c, err := m.MangaDetailsScrapper.GetContent(m.buildPath(mangaId))
	if err != nil {
		return contract.Manga{}, err
	}

	info := contract.MangaInfo{
		Author:      c.Author,
		Status:      c.Status,
		ReleaseYear: c.ReleaseYear,
		Description: c.Description,
		Genres:      c.Genres,
	}

	var chapters []contract.MangaChapter
	for _, chapter := range c.Chapters {
		chapters = append(chapters, contract.MangaChapter{
			ID:           chapter.ID,
			Title:        chapter.Title,
			LastModified: chapter.LastModified,
		})
	}

	return contract.Manga{
		PageTitle:      c.Name,
		HeaderImageURL: c.IconURL,
		Info:           info,
		Chapters:       chapters,
	}, nil
}

func (m mangasailService) GetChapterDetails(chapterId string) (contract.Chapter, error) {
	c, err := m.ChapterDetailsScrapper.GetContent(m.buildPath(chapterId))
	if err != nil {
		return contract.Chapter{}, err
	}

	var imageURLs []string
	for _, image := range c.Images {
		imageURLs = append(imageURLs, image.ImageURL)
	}

	return contract.Chapter{
		PageTitle: c.Title,
		ImageURLs: imageURLs,
	}, nil
}

func (m mangasailService) buildDailyHotMangasHomeCard(mangas domain.Mangas) contract.HomeCard {
	card := contract.HomeCard{
		Identifier: contract.HomeCardDailyHotMangas,
		Label:      "Daily Hot Manga Chapter",
		Content:    []contract.HomeCardContent{},
	}
	for _, manga := range mangas {
		if len(manga.Chapters) > 0 {
			card.Content = append(
				card.Content,
				contract.HomeCardContent{
					Title:   manga.Chapters[0].Title,
					IconURL: manga.IconURL,
					Redirection: contract.HomeCardRedirection{
						ID:   manga.Chapters[0].ID,
						Type: contract.RedirectToChapter,
					},
				},
			)
		}
	}
	return card
}

func (m mangasailService) buildPopularMangasHomeCard(mangas domain.Mangas) contract.HomeCard {
	card := contract.HomeCard{
		Identifier: contract.HomeCardPopularMangas,
		Label:      "Popular Manga",
		Content:    []contract.HomeCardContent{},
	}
	for _, manga := range mangas {
		if len(manga.ID) > 0 {
			card.Content = append(
				card.Content,
				contract.HomeCardContent{
					Title:   manga.Name,
					IconURL: manga.IconURL,
					Redirection: contract.HomeCardRedirection{
						ID:   manga.ID,
						Type: contract.RedirectToManga,
					},
				},
			)
		}
	}
	return card
}

func (m mangasailService) buildLatestMangasHomeCard(mangas domain.Mangas) contract.HomeCard {
	card := contract.HomeCard{
		Identifier: contract.HomeCardLatestMangas,
		Label:      "Latest Updated Manga",
		Content:    []contract.HomeCardContent{},
	}
	for _, manga := range mangas {
		if len(manga.Chapters) > 0 {
			card.Content = append(
				card.Content,
				contract.HomeCardContent{
					Title:    manga.Name,
					SubTitle: manga.Chapters[0].Title,
					Info:     manga.Chapters[0].LastModified,
					IconURL:  manga.IconURL,
					Redirection: contract.HomeCardRedirection{
						ID:   manga.ID,
						Type: contract.RedirectToManga,
					},
				},
			)
		}
	}
	return card
}

func (m mangasailService) buildNewMangasHomeCard(mangas domain.Mangas) contract.HomeCard {
	card := contract.HomeCard{
		Identifier: contract.HomeCardNewMangas,
		Label:      "New Released Manga",
		Content:    []contract.HomeCardContent{},
	}
	for _, manga := range mangas {
		if len(manga.ID) > 0 {
			card.Content = append(
				card.Content,
				contract.HomeCardContent{
					Title:   manga.Name,
					IconURL: manga.IconURL,
					Redirection: contract.HomeCardRedirection{
						ID:   manga.ID,
						Type: contract.RedirectToManga,
					},
				},
			)

		}
	}
	return card
}

func (m mangasailService) buildPath(id string) string {
	return common.MangasailPrefixPath + id
}
