package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/service"
)

func getSourceFromRequest(r *http.Request) (string, error) {
	return getValueFromRequest(r, common.ParamKeySource, "unknown source")
}

func getIdFromRequest(r *http.Request) (string, error) {
	return getValueFromRequest(r, common.ParamKeyId, "unknown id")
}

func getValueFromRequest(r *http.Request, key, errTxt string) (string, error) {
	s := r.FormValue(key)
	if len(s) == 0 {
		return "", errors.New(errTxt)
	}
	return s, nil
}

func getFetchServiceFromRequest(r *http.Request, sMap map[common.FetchServiceKey]service.FetchService) (service.FetchService, error) {
	source, err := getSourceFromRequest(r)
	if err != nil {
		return nil, err
	}

	svc, ok := sMap[common.FetchServiceKey(source)]
	if !ok {
		return nil, errors.New("unknown source")
	}

	return svc, nil
}

func respondWith(statusCode int, w http.ResponseWriter, response contract.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
