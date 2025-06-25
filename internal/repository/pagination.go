package repository

import (
	"net/http"
	"strconv"
)

type PaginatedVillaQuery struct {
	Limit    int    `json:"limit" validate:"gte=1,lte=10"`
	Offset   int    `json:"offset" validate:"gte=0"`
	Sort     string `json:"sort" validate:"oneof=asc desc"`
	Location string `json:"location"`
	Category string `json:"category"`
	MinGuest string `json:"min_guest"`
	Bedrooms string `json:"bedrooms"`
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
	if sort != "" {
		pv.Sort = sort
	}

	location := qs.Get("location")
	if location != "" {
		pv.Location = location
	}

	category := qs.Get("category")
	if category != "" {
		pv.Category = category
	}

	bedrooms := qs.Get("bedrooms")
	if bedrooms != "" {
		pv.Bedrooms = bedrooms
	}

	minGuest := qs.Get("min_guest")
	if minGuest != "" {
		pv.MinGuest = minGuest
	}

	return pv, nil
}
