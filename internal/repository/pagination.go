package repository

import (
	"net/http"
	"strconv"
)

type PaginatedVillaQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pv PaginatedVillaQuery) Parse(r *http.Request) (PaginatedVillaQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pv, err
		}

		pv.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return pv, err
		}

		pv.Offset = of
	}

	sort := qs.Get("sort")
	if offset != "" {
		pv.Sort = sort
	}

	return pv, nil
}
