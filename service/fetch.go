package service

import "github.com/bigscreen/manga-scrapper/contract"

type FetchService interface {
	GetHomeCards() (contract.Home, error)
	GetMangaDetails(mangaId string) (contract.Manga, error)
	GetChapterDetails(chapterId string) (contract.Chapter, error)
}
