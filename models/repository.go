package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductList struct {
	Products      []Product
	TotalProducts int
}

type GetProductsFilter struct {
	Limit    int
	Offset   int
	Category string
	Price    float64
	Page     int
}

type DataStore interface {
	GetAllProducts(query *GetProductsFilter) (*ProductList, error)
	GetProduct(code string) (*Product, error)
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts(query *GetProductsFilter) (*ProductList, error) {
	var products []Product
	var total int64

	q := r.db.Model(&Product{})
	if query.Category != "" {
		q = q.Joins("JOIN categories ON categories.product_id = products.id").
			Where("categories.name = ?", query.Category)
	}

	if query.Price > 0 {
		q = q.Where("products.price < ?", query.Price)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := q.
		Limit(query.Limit).Offset(query.Offset).Preload("Variants").
		Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return &ProductList{Products: products, TotalProducts: int(total)}, nil
}

func (r *ProductsRepository) GetProduct(code string) (*Product, error) {
	var product Product

	if err := r.db.Preload("Variants").
		Preload("Category").
		Where("code = ?", code).
		First(&product).Error; err != nil {
		return nil, err
	}

	for i := range product.Variants {
		if decimal.Zero.Equal(product.Variants[i].Price) {
			product.Variants[i].Price = product.Price
		}
	}
	return &product, nil
}
