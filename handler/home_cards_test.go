package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/config"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/logger"
	"github.com/bigscreen/manga-scrapper/mock"
	"github.com/bigscreen/manga-scrapper/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

func TestHomeCardsHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HomeCardsHandlerTestSuite))
}

type HomeCardsHandlerTestSuite struct {
	suite.Suite
	mFetchService *mock.FetchServiceMock
	mr            *mux.Router
}

func (h *HomeCardsHandlerTestSuite) SetupSuite() {
	config.Load()
	logger.SetupLogger()
}

func (h *HomeCardsHandlerTestSuite) SetupTest() {
	h.mFetchService = &mock.FetchServiceMock{}
	h.mr = mux.NewRouter()
}

func (h *HomeCardsHandlerTestSuite) TestGetHomeCards() {
	cases := []struct {
		name             string
		source           string
		serviceResult    *contract.Home
		serviceError     error
		expectedResponse contract.Response
		expectedRespCode int
	}{
		{
			name:             "WhenSourceIsEmpty",
			expectedResponse: contract.NewErrorResponse(errors.New("unknown source")),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenSourceIsUnknown",
			source:           "foo",
			expectedResponse: contract.NewErrorResponse(errors.New("unknown source")),
			expectedRespCode: http.StatusBadRequest,
		},
		{
			name:             "WhenSourceIsKnownAndSuccessFromService",
			source:           string(common.FSKeyMangasail),
			serviceResult:    &contract.Home{PageTitle: "Home"},
			serviceError:     nil,
			expectedResponse: contract.NewSuccessResponse(contract.Home{PageTitle: "Home"}),
			expectedRespCode: http.StatusOK,
		},
		{
			name:             "WhenSourceIsKnownAndErrorFromService",
			source:           string(common.FSKeyMangasail),
			serviceResult:    &contract.Home{},
			serviceError:     errors.New("some error"),
			expectedResponse: contract.NewErrorResponse(errors.New("some error")),
			expectedRespCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		h.SetupTest()
		h.Run(tc.name, func() {
			if tc.serviceResult != nil || tc.serviceError != nil {
				h.mFetchService.On("GetHomeCards").Return(*tc.serviceResult, tc.serviceError)
			}
			req, rr := h.buildRequest(tc.source)
			svcMap := map[common.FetchServiceKey]service.FetchService{common.FSKeyMangasail: h.mFetchService}
			h.mr.HandleFunc(common.GetHomeCardsAPIPath, GetHomeCards(svcMap))
			h.mr.ServeHTTP(rr, req)

			h.Equal(tc.expectedResponse.String(), strings.TrimSuffix(rr.Body.String(), "\n"))
			h.Equal(tc.expectedRespCode, rr.Code)

			h.mFetchService.AssertExpectations(h.T())
		})
	}
}

func (h *HomeCardsHandlerTestSuite) buildRequest(source string) (*http.Request, *httptest.ResponseRecorder) {
	apiPath := fmt.Sprintf("%s?source=%s", common.GetHomeCardsAPIPath, source)
	req, _ := http.NewRequest("GET", apiPath, nil)
	rr := httptest.NewRecorder()
	return req, rr
}
