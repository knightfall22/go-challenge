package product

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductHandler struct {
	repo models.DataStore
}

type Response struct {
	Product Product `json:"product"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Variant  []Variant `json:"variant"`
}

type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

func NewCatalogHandler(r models.DataStore) *ProductHandler {
	return &ProductHandler{
		repo: r,
	}
}

func (h *ProductHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	product, err := h.repo.GetProduct(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	// Map response
	var variants []Variant
	for _, v := range product.Variants {
		variants = append(variants, Variant{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: v.Price.InexactFloat64(),
		})
	}

	response := Response{
		Product{
			Code:     product.Code,
			Price:    product.Price.InexactFloat64(),
			Category: product.Category.Name,
			Variant:  variants,
		},
	}
	api.OKResponse(w, http.StatusOK, api.ApiResponse[Response]{
		Success: true,
		Message: "Product found",
		Data:    response,
	})

}
