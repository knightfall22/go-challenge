package catalog

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

type CatalogHandler struct {
	repo models.DataStore
}

func NewCatalogHandler(r models.DataStore) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	filter, err := h.parseCatalogFilters(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.repo.GetAllProducts(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Map response
	products := make([]Product, len(res.Products))
	for i, p := range res.Products {
		products[i] = Product{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: p.Category.Name,
		}
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Products: products,
	}

	api.OKResponse(
		w,
		http.StatusOK,
		api.ApiResponse[Response]{
			Success: true,
			Message: "Fetch Catalog successfully",
			Data:    response,

			Limit:   filter.Limit,
			Page:    filter.Page,
			Count:   len(products),
			HasNext: (filter.Page * filter.Limit) < res.TotalProducts,
			Total:   res.TotalProducts,
		},
	)

}

func (h *CatalogHandler) parseCatalogFilters(r *http.Request) (*models.GetProductsFilter, error) {
	query := r.URL.Query()

	limit := 10
	page := 1

	if limitStr := query.Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'limit' parameter")
		}
		limit = parsedLimit
	}

	if limit > 100 || limit < 1 {
		//reset limit to default when it exceeds this bound
		limit = 10
	}

	if pageStr := query.Get("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid 'limit' parameter")
		}
		page = parsedPage
	}

	priceLessThan := 0.0
	if priceStr := query.Get("price"); priceStr != "" {
		parsedprice, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid 'limit' parameter")
		}
		priceLessThan = parsedprice
	}

	offest := (page - 1) * limit

	return &models.GetProductsFilter{
		Limit:    limit,
		Page:     page,
		Offset:   offest,
		Category: query.Get("category"),
		Price:    priceLessThan,
	}, nil
}
