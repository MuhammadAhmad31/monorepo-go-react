package mapper

import (
	"backend/internal/generated"
	"backend/internal/models"
)

func ToGeneratedProduct(product *models.Product) generated.Product {
	return generated.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   &product.CreatedAt,
		UpdatedAt:   &product.UpdatedAt,
	}
}

func ToGeneratedProducts(products []models.Product) []generated.Product {
	result := make([]generated.Product, len(products))
	for i := range products {
		result[i] = ToGeneratedProduct(&products[i])
	}
	return result
}
