package mock

import (
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/stretchr/testify/mock"
)

type FetchServiceMock struct {
	mock.Mock
}

func (f *FetchServiceMock) GetHomeCards() (contract.Home, error) {
	args := f.Called()
	if args.Get(1) != nil {
		return contract.Home{}, args.Get(1).(error)
	}
	return args.Get(0).(contract.Home), nil
}

func (f *FetchServiceMock) GetMangaDetails(mangaId string) (contract.Manga, error) {
	args := f.Called(mangaId)
	if args.Get(1) != nil {
		return contract.Manga{}, args.Get(1).(error)
	}
	return args.Get(0).(contract.Manga), nil
}

func (f *FetchServiceMock) GetChapterDetails(chapterId string) (contract.Chapter, error) {
	args := f.Called(chapterId)
	if args.Get(1) != nil {
		return contract.Chapter{}, args.Get(1).(error)
	}
	return args.Get(0).(contract.Chapter), nil
}
