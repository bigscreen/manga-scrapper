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

func TestMangaDetailsHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MangaDetailsHandlerTestSuite))
}

type MangaDetailsHandlerTestSuite struct {
	suite.Suite
	mFetchService *mock.FetchServiceMock
	mr            *mux.Router
}

func (m *MangaDetailsHandlerTestSuite) SetupSuite() {
	config.Load()
	logger.SetupLogger()
}

func (m *MangaDetailsHandlerTestSuite) SetupTest() {
	m.mFetchService = &mock.FetchServiceMock{}
	m.mr = mux.NewRouter()
}

func (m *MangaDetailsHandlerTestSuite) TestGetMangaDetails() {
	cases := []struct {
		name             string
		mangaId          string
		source           string
		serviceResult    *contract.Manga
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
			mangaId:          "mangaID",
			expectedResponse: contract.NewErrorResponse(errors.New(errors.WithValidationErrorCode(), errors.WithMessage("Source is missing."))),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenIdExistsAndSourceIsUnknown",
			mangaId:          "mangaID",
			source:           "foo",
			expectedResponse: contract.NewErrorResponse(errors.New(errors.WithValidationErrorCode(), errors.WithMessage("Unknown source."))),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenIdExistsAndSourceIsKnownAndSuccessFromService",
			mangaId:          "mangaID",
			source:           string(common.FSKeyMangasail),
			serviceResult:    &contract.Manga{PageTitle: "MangaA"},
			serviceError:     nil,
			expectedResponse: contract.NewSuccessResponse(contract.Manga{PageTitle: "MangaA"}),
			expectedRespCode: http.StatusOK,
		},
		{
			name:             "WhenIdExistsAndSourceIsKnownAndErrorFromService",
			mangaId:          "mangaID",
			source:           string(common.FSKeyMangasail),
			serviceResult:    &contract.Manga{},
			serviceError:     nativeErrs.New("some error"),
			expectedResponse: contract.NewErrorResponse(nativeErrs.New("some error")),
			expectedRespCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		m.SetupTest()
		m.Run(tc.name, func() {
			if tc.serviceResult != nil || tc.serviceError != nil {
				m.mFetchService.On("GetMangaDetails", tc.mangaId).Return(*tc.serviceResult, tc.serviceError)
			}
			req, rr := m.buildRequest(tc.mangaId, tc.source)
			svcMap := map[common.FetchServiceKey]service.FetchService{common.FSKeyMangasail: m.mFetchService}
			m.mr.HandleFunc(common.GetMangaDetailsAPIPath, GetMangaDetails(svcMap))
			m.mr.ServeHTTP(rr, req)

			m.Equal(tc.expectedResponse.String(), strings.TrimSuffix(rr.Body.String(), "\n"))
			m.Equal(tc.expectedRespCode, rr.Code)

			m.mFetchService.AssertExpectations(m.T())
		})
	}
}

func (m *MangaDetailsHandlerTestSuite) buildRequest(id string, source string) (*http.Request, *httptest.ResponseRecorder) {
	apiPath := fmt.Sprintf("%s?id=%s&source=%s", common.GetMangaDetailsAPIPath, id, source)
	req, _ := http.NewRequest("GET", apiPath, nil)
	rr := httptest.NewRecorder()
	return req, rr
}
