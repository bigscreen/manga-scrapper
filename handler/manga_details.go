package handler

import (
	"net/http"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/service"
)

func GetMangaDetails(svcMap map[common.FetchServiceKey]service.FetchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromRequest(r)
		if err != nil {
			responseError(w, err)
			return
		}

		svc, err := getFetchServiceFromRequest(r, svcMap)
		if err != nil {
			responseError(w, err)
			return
		}

		mr, err := svc.GetMangaDetails(id)
		if err != nil {
			responseError(w, err)
			return
		}
		respondWith(http.StatusOK, w, contract.NewSuccessResponse(mr))
	}
}
