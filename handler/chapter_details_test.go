package handler

import (
	nativeErrs "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/config"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/errors"
	"github.com/bigscreen/manga-scrapper/logger"
	"github.com/bigscreen/manga-scrapper/mock"
	"github.com/bigscreen/manga-scrapper/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

func TestChapterDetailsHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ChapterDetailsHandlerTestSuite))
}

type ChapterDetailsHandlerTestSuite struct {
	suite.Suite
	mFetchService *mock.FetchServiceMock
	mr            *mux.Router
}

func (c *ChapterDetailsHandlerTestSuite) SetupSuite() {
	config.Load()
	logger.SetupLogger()
}

func (c *ChapterDetailsHandlerTestSuite) SetupTest() {
	c.mFetchService = &mock.FetchServiceMock{}
	c.mr = mux.NewRouter()
}

func (c *ChapterDetailsHandlerTestSuite) TestGetMangaDetails() {
	cases := []struct {
		name             string
		chapterId        string
		source           string
		serviceResult    *contract.Chapter
		serviceError     error
		expectedResponse contract.Response
		expectedRespCode int
	}{
		{
			name:             "WhenIdIsMissing",
			expectedResponse: contract.NewErrorResponse(errors.New(errors.WithValidationErrorCode(), errors.WithMessage("Id is missing."))),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenIdExistsAndSourceIsEmpty",
			chapterId:        "chapterID",
			expectedResponse: contract.NewErrorResponse(errors.New(errors.WithValidationErrorCode(), errors.WithMessage("Source is missing."))),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenIdExistsAndSourceIsUnknown",
			chapterId:        "chapterID",
			source:           "foo",
			expectedResponse: contract.NewErrorResponse(errors.New(errors.WithValidationErrorCode(), errors.WithMessage("Unknown source."))),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenIdExistsAndSourceIsKnownAndSuccessFromService",
			chapterId:        "chapterID",
			source:           string(common.FSKeyMangasail),
			serviceResult:    &contract.Chapter{PageTitle: "ChapterA"},
			serviceError:     nil,
			expectedResponse: contract.NewSuccessResponse(contract.Chapter{PageTitle: "ChapterA"}),
			expectedRespCode: http.StatusOK,
		},
		{
			name:             "WhenIdExistsAndSourceIsKnownAndErrorFromService",
			chapterId:        "chapterID",
			source:           string(common.FSKeyMangasail),
			serviceResult:    &contract.Chapter{},
			serviceError:     nativeErrs.New("some error"),
			expectedResponse: contract.NewErrorResponse(nativeErrs.New("some error")),
			expectedRespCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		c.SetupTest()
		c.Run(tc.name, func() {
			if tc.serviceResult != nil || tc.serviceError != nil {
				c.mFetchService.On("GetChapterDetails", tc.chapterId).Return(*tc.serviceResult, tc.serviceError)
			}
			req, rr := c.buildRequest(tc.chapterId, tc.source)
			svcMap := map[common.FetchServiceKey]service.FetchService{common.FSKeyMangasail: c.mFetchService}
			c.mr.HandleFunc(common.GetChapterDetailsAPIPath, GetChapterDetails(svcMap))
			c.mr.ServeHTTP(rr, req)

			c.Equal(tc.expectedResponse.String(), strings.TrimSuffix(rr.Body.String(), "\n"))
			c.Equal(tc.expectedRespCode, rr.Code)

			c.mFetchService.AssertExpectations(c.T())
		})
	}
}

func (c *ChapterDetailsHandlerTestSuite) buildRequest(id string, source string) (*http.Request, *httptest.ResponseRecorder) {
	apiPath := fmt.Sprintf("%s?id=%s&source=%s", common.GetChapterDetailsAPIPath, id, source)
	req, _ := http.NewRequest("GET", apiPath, nil)
	rr := httptest.NewRecorder()
	return req, rr
}
