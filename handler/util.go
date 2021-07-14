package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bigscreen/manga-scrapper/common"
	"github.com/bigscreen/manga-scrapper/contract"
	"github.com/bigscreen/manga-scrapper/errors"
	"github.com/bigscreen/manga-scrapper/service"
)

func getSourceFromRequest(r *http.Request) (string, error) {
	return getValueFromRequest(r, common.ParamKeySource, "Source is missing.")
}

func getIdFromRequest(r *http.Request) (string, error) {
	return getValueFromRequest(r, common.ParamKeyId, "Id is missing.")
}

func getValueFromRequest(r *http.Request, key, errMsg string) (string, error) {
	s := r.FormValue(key)
	if len(s) == 0 {
		return "", errors.New(errors.WithValidationErrorCode(), errors.WithMessage(errMsg))
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
		return nil, errors.New(errors.WithValidationErrorCode(), errors.WithMessage("Unknown source."))
	}

	return svc, nil
}

func respondWith(statusCode int, w http.ResponseWriter, response contract.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func responseError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if e := errors.FromError(err); e != nil {
		code = e.HttpStatusCode()
	}
	respondWith(code, w, contract.NewErrorResponse(err))
}
