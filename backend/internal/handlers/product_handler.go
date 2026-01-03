package handlers

import (
	"backend/internal/generated"
	"backend/internal/handlers/mapper"
	"backend/internal/models"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ListProducts(c *gin.Context, params generated.ListProductsParams) {
	page := 1
	perPage := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.PerPage != nil {
		perPage = *params.PerPage
	}

	products, total, err := h.service.ListProducts(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to fetch products",
		})
		return
	}

	totalInt := int(total)

	c.JSON(http.StatusOK, gin.H{
		"data": mapper.ToGeneratedProducts(products),
		"meta": generated.Meta{
			Page:    &page,
			PerPage: &perPage,
			Total:   &totalInt,
		},
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req generated.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request body",
		})
		return
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       req.Stock,
		Category:    req.Category,
	}

	if err := h.service.CreateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": mapper.ToGeneratedProduct(product),
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context, id generated.IdParam) {
	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, generated.Error{
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mapper.ToGeneratedProduct(product),
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context, id generated.IdParam) {
	var req generated.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request body",
		})
		return
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       req.Stock,
		Category:    req.Category,
	}

	if err := h.service.UpdateProduct(c.Request.Context(), id, product); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, generated.Error{
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mapper.ToGeneratedProduct(product),
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context, id generated.IdParam) {
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, generated.Error{
				Message: "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to delete product",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
