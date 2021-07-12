package handler

import (
	"net/http"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/service"
)

func GetHomeCards(svcMap map[common.FetchServiceKey]service.FetchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc, err := getFetchServiceFromRequest(r, svcMap)
		if err != nil {
			respondWith(http.StatusBadRequest, w, contract.NewErrorResponse(err))
			return
		}

		hr, err := svc.GetHomeCards()
		if err != nil {
			respondWith(http.StatusInternalServerError, w, contract.NewErrorResponse(err))
			return
		}
		respondWith(http.StatusOK, w, contract.NewSuccessResponse(hr))
	}
}
