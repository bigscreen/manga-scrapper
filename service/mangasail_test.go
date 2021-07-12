package service

import (
	"errors"
	"testing"

	"github.com/bigscreen/manga-scrapper/config"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/domain"
	"github.com/bigscreen/manga-scrapper/logger"
	"github.com/bigscreen/manga-scrapper/mock"
	"github.com/stretchr/testify/suite"
)

func TestMangasailServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MangasailServiceTestSuite))
}

type MangasailServiceTestSuite struct {
	suite.Suite
	mHomeScrapper           *mock.MangasailHomeScrapperMock
	mMangaDetailsScrapper   *mock.MangasailMangaDetailsScrapperMock
	mChapterDetailsScrapper *mock.MangasailChapterDetailsScrapperMock
	mangasailSvc            FetchService
}

func (m *MangasailServiceTestSuite) SetupSuite() {
	config.Load()
	logger.SetupLogger()
}

func (m *MangasailServiceTestSuite) SetupTest() {
	m.mHomeScrapper = &mock.MangasailHomeScrapperMock{}
	m.mMangaDetailsScrapper = &mock.MangasailMangaDetailsScrapperMock{}
	m.mChapterDetailsScrapper = &mock.MangasailChapterDetailsScrapperMock{}
	m.mangasailSvc = NewMangsailService(MangasailServiceParams{
		HomeScrapper:           m.mHomeScrapper,
		MangaDetailsScrapper:   m.mMangaDetailsScrapper,
		ChapterDetailsScrapper: m.mChapterDetailsScrapper,
	})
}

func (m *MangasailServiceTestSuite) TestGetHomeCards() {
	cases := []struct {
		name           string
		scrapperResult domain.HomeMangas
		scrapperError  error
		expectedResult contract.Home
		expectedError  error
	}{
		{
			name: "WhenSuccessFromScrapper",
			scrapperResult: domain.HomeMangas{
				DailyHotMangas: m.getMangas(),
				LatestMangas:   m.getMangas(),
				PopularMangas:  m.getMangas(),
				NewMangas:      m.getMangas(),
			},
			expectedResult: contract.Home{
				PageTitle: "Mangajack Home",
				Cards: []contract.HomeCard{
					{
						Identifier: contract.HomeCardDailyHotMangas,
						Label:      "Daily Hot Manga Chapter",
						Content: []contract.HomeCardContent{
							{
								Title:   "Manga A 76",
								IconURL: "https://icon.jpg",
								Redirection: contract.HomeCardRedirection{
									ID:   "chapterID",
									Type: contract.RedirectToChapter,
								},
							},
						},
					},
					{
						Identifier: contract.HomeCardPopularMangas,
						Label:      "Popular Manga",
						Content: []contract.HomeCardContent{
							{
								Title:   "Manga A",
								IconURL: "https://icon.jpg",
								Redirection: contract.HomeCardRedirection{
									ID:   "mangaID",
									Type: contract.RedirectToManga,
								},
							},
						},
					},
					{
						Identifier: contract.HomeCardLatestMangas,
						Label:      "Latest Updated Manga",
						Content: []contract.HomeCardContent{
							{
								Title:    "Manga A",
								SubTitle: "Manga A 76",
								Info:     "12 Jul 2021",
								IconURL:  "https://icon.jpg",
								Redirection: contract.HomeCardRedirection{
									ID:   "mangaID",
									Type: contract.RedirectToManga,
								},
							},
						},
					},
					{
						Identifier: contract.HomeCardNewMangas,
						Label:      "New Released Manga",
						Content: []contract.HomeCardContent{
							{
								Title:   "Manga A",
								IconURL: "https://icon.jpg",
								Redirection: contract.HomeCardRedirection{
									ID:   "mangaID",
									Type: contract.RedirectToManga,
								},
							},
						},
					},
				},
			},
		},
		{
			name:          "WhenErrorFromScrapper",
			scrapperError: errors.New("foo"),
			expectedError: errors.New("foo"),
		},
	}

	for _, tc := range cases {
		m.SetupTest()
		m.Suite.Run(tc.name, func() {
			m.mHomeScrapper.On("GetContent").Return(tc.scrapperResult, tc.scrapperError)

			result, err := m.mangasailSvc.GetHomeCards()

			m.Equal(tc.expectedResult, result)
			m.Equal(tc.expectedError, err)
			m.mHomeScrapper.AssertExpectations(m.T())
		})
	}
}

func (m *MangasailServiceTestSuite) TestGetMangaDetails() {
	cases := []struct {
		name           string
		scrapperResult domain.Manga
		scrapperError  error
		expectedResult contract.Manga
		expectedError  error
	}{
		{
			name:           "WhenSuccessFromScrapper",
			scrapperResult: m.getMangas()[0],
			expectedResult: contract.Manga{
				PageTitle:      "Manga A",
				HeaderImageURL: "https://icon.jpg",
				Info: contract.MangaInfo{
					Author:      "Mr X",
					Status:      "Ongoing",
					ReleaseYear: "2021",
					Description: "Lorem ipsum doler sit amet",
					Genres:      []string{"action"},
				},
				Chapters: []contract.MangaChapter{
					{
						ID:           "chapterID",
						Title:        "Manga A 76",
						LastModified: "12 Jul 2021",
					},
				},
			},
		},
		{
			name:          "WhenErrorFromScrapper",
			scrapperError: errors.New("foo"),
			expectedError: errors.New("foo"),
		},
	}

	for _, tc := range cases {
		m.SetupTest()
		m.Suite.Run(tc.name, func() {
			m.mMangaDetailsScrapper.On("GetContent", "/content/mangaID").
				Return(tc.scrapperResult, tc.scrapperError)

			result, err := m.mangasailSvc.GetMangaDetails("mangaID")

			m.Equal(tc.expectedResult, result)
			m.Equal(tc.expectedError, err)
			m.mMangaDetailsScrapper.AssertExpectations(m.T())
		})
	}
}

func (m *MangasailServiceTestSuite) TestGetChapterDetails() {
	cases := []struct {
		name           string
		scrapperResult domain.Chapter
		scrapperError  error
		expectedResult contract.Chapter
		expectedError  error
	}{
		{
			name:           "WhenSuccessFromScrapper",
			scrapperResult: m.getChapters()[0],
			expectedResult: contract.Chapter{
				PageTitle: "Manga A 76",
				ImageURLs: []string{"https://image.jpg"},
			},
		},
		{
			name:          "WhenErrorFromScrapper",
			scrapperError: errors.New("foo"),
			expectedError: errors.New("foo"),
		},
	}

	for _, tc := range cases {
		m.SetupTest()
		m.Suite.Run(tc.name, func() {
			m.mChapterDetailsScrapper.On("GetContent", "/content/chapterID").
				Return(tc.scrapperResult, tc.scrapperError)

			result, err := m.mangasailSvc.GetChapterDetails("chapterID")

			m.Equal(tc.expectedResult, result)
			m.Equal(tc.expectedError, err)
			m.mChapterDetailsScrapper.AssertExpectations(m.T())
		})
	}
}

func (m *MangasailServiceTestSuite) getMangas() domain.Mangas {
	return domain.Mangas{
		domain.Manga{
			ID:          "mangaID",
			Name:        "Manga A",
			IconURL:     "https://icon.jpg",
			Author:      "Mr X",
			Status:      "Ongoing",
			ReleaseYear: "2021",
			Description: "Lorem ipsum doler sit amet",
			Genres:      []string{"action"},
			Chapters:    m.getChapters(),
		},
	}
}

func (m *MangasailServiceTestSuite) getChapters() domain.Chapters {
	return domain.Chapters{
		domain.Chapter{
			ID:           "chapterID",
			Title:        "Manga A 76",
			LastModified: "12 Jul 2021",
			Images:       m.getImages(),
		},
	}
}

func (m *MangasailServiceTestSuite) getImages() domain.Images {
	return domain.Images{
		domain.Image{
			ID:       "1",
			ImageURL: "https://image.jpg",
		},
	}
}
