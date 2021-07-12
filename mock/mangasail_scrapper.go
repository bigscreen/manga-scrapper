package mock

import (
	"github.com/bigscreen/manga-scrapper/domain"
	"github.com/stretchr/testify/mock"
)

type MangasailHomeScrapperMock struct {
	mock.Mock
}

func (m *MangasailHomeScrapperMock) GetContent() (domain.HomeMangas, error) {
	args := m.Called()
	if args.Get(1) != nil {
		return domain.HomeMangas{}, args.Get(1).(error)
	}
	return args.Get(0).(domain.HomeMangas), nil
}

type MangasailMangaDetailsScrapperMock struct {
	mock.Mock
}

func (m *MangasailMangaDetailsScrapperMock) GetContent(path string) (domain.Manga, error) {
	args := m.Called(path)
	if args.Get(1) != nil {
		return domain.Manga{}, args.Get(1).(error)
	}
	return args.Get(0).(domain.Manga), nil
}

type MangasailChapterDetailsScrapperMock struct {
	mock.Mock
}

func (m *MangasailChapterDetailsScrapperMock) GetContent(path string) (domain.Chapter, error) {
	args := m.Called(path)
	if args.Get(1) != nil {
		return domain.Chapter{}, args.Get(1).(error)
	}
	return args.Get(0).(domain.Chapter), nil
}
