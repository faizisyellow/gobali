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

type PaginatedLocationQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pl PaginatedLocationQuery) Parse(r *http.Request) (PaginatedLocationQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pl, err
		}

		pl.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return pl, err
		}

		pl.Offset = of
	}

	sort := qs.Get("sort")
	if sort != "" {
		pl.Sort = sort
	}

	return pl, nil
}

type PaginatedCategoriesQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pc PaginatedCategoriesQuery) Parse(r *http.Request) (PaginatedCategoriesQuery, error) {

	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pc, err
		}

		pc.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return pc, err
		}

		pc.Offset = of
	}

	sort := qs.Get("sort")
	if sort != "" {
		pc.Sort = sort
	}

	return pc, nil
}

type PaginatedAmenitiesQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pa PaginatedAmenitiesQuery) Parse(r *http.Request) (PaginatedAmenitiesQuery, error) {

	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pa, err
		}

		pa.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return pa, err
		}

		pa.Offset = of
	}

	sort := qs.Get("sort")
	if sort != "" {
		pa.Sort = sort
	}

	return pa, nil
}

type PaginatedBookingsQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pb PaginatedBookingsQuery) Parse(r *http.Request) (PaginatedBookingsQuery, error) {

	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pb, err
		}

		pb.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return pb, err
		}

		pb.Offset = of
	}

	sort := qs.Get("sort")
	if sort != "" {
		pb.Sort = sort
	}

	return pb, nil
}

type PaginatedUserBookingsQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=10"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (pb PaginatedUserBookingsQuery) Parse(r *http.Request) (PaginatedUserBookingsQuery, error) {

	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return pb, err
		}

		pb.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		of, err := strconv.Atoi(offset)
		if err != nil {
			return pb, err
		}

		pb.Offset = of
	}

	sort := qs.Get("sort")
	if sort != "" {
		pb.Sort = sort
	}

	return pb, nil
}
