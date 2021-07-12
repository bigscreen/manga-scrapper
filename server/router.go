package server

import (
	"net/http"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/dep"
	"github.com/bigscreen/manga-scrapper/handler"
	"github.com/gorilla/mux"
)

func Router(deps dep.AppDependencies) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", handler.PingHandler).Methods("GET")
	router.HandleFunc(common.GetHomeCardsAPIPath, handler.GetHomeCards(deps.FetchServiceMap)).Methods("GET")
	router.HandleFunc(common.GetMangaDetailsAPIPath, handler.GetMangaDetails(deps.FetchServiceMap)).Methods("GET")
	router.HandleFunc(common.GetChapterDetailsAPIPath, handler.GetChapterDetails(deps.FetchServiceMap)).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)

	return router
}
