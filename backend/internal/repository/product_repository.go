package repository

import (
	"backend/internal/generated"
	"backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	FindByID(ctx context.Context, id generated.IdParam) (*models.Product, error)
	FindAll(ctx context.Context, page, perPage int) ([]models.Product, int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id generated.IdParam) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, id generated.IdParam) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(ctx context.Context, page, perPage int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	offset := (page - 1) * perPage

	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(perPage).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id generated.IdParam) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}
