package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"
	"faizalmaulana/lsp/models/entity"

	"github.com/gin-gonic/gin"
)

type ItemsHandler struct {
	cfg    *conf.Config
	items  services.ItemsService
	images services.ImagesService
}

func NewItemsHandler(cfg *conf.Config, items services.ItemsService, images services.ImagesService) *ItemsHandler {
	return &ItemsHandler{cfg: cfg, items: items, images: images}
}

func (h *ItemsHandler) Register(rr *gin.RouterGroup) {
	rg := rr.Group("/items")
	rg.GET("", h.list)
	rg.GET(":id", h.get)
	rg.POST("", middleware.JWTMiddleware(h.cfg), h.create)
	rg.PUT(":id", middleware.JWTMiddleware(h.cfg), h.update)
	rg.DELETE(":id", middleware.JWTMiddleware(h.cfg), h.delete)
}

func (h *ItemsHandler) list(c *gin.Context) {
	count := 10
	page := 1
	if v := c.Query("count"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			count = n
		}
	}
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			page = n
		}
	}
	items, err := h.items.GetAll(count, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to list items"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", items))
}

func (h *ItemsHandler) get(c *gin.Context) {
	id := c.Param("id")
	item, err := h.items.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("item not found"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", item))
}

func (h *ItemsHandler) create(c *gin.Context) {
	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}

	fmt.Println(c.Request.RequestURI)
	imageFileName := req.ImageUrl
	if h.images != nil && req.ImageBase64 != "" {
		_, stored, err := h.images.UploadBase64(req.ItemName, req.ImageType, req.ImageBase64)
		if err == nil {
			imageFileName = stored
		}
	}
	it := &entity.Items{IdItem: helper.Uuid(), ItemName: req.ItemName, ItemType: req.ItemType, Price: req.Price, Description: req.Description, ImageUrl: imageFileName}
	if req.IsAvailable != nil {
		it.IsAvailable = *req.IsAvailable
	}
	saved, err := h.items.Create(it)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to create item"))
		return
	}
	c.JSON(http.StatusCreated, helper.SuccessResponse("created", saved))
}

func (h *ItemsHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	existing, err := h.items.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("item not found"))
		return
	}
	if req.ItemName != nil {
		existing.ItemName = *req.ItemName
	}
	if req.ItemType != nil {
		existing.ItemType = *req.ItemType
	}
	if req.IsAvailable != nil {
		existing.IsAvailable = *req.IsAvailable
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.ImageUrl != nil {
		existing.ImageUrl = *req.ImageUrl
	}
	if h.images != nil && req.ImageBase64 != nil && *req.ImageBase64 != "" {
		ct := ""
		if req.ImageType != nil {
			ct = *req.ImageType
		}
		_, stored, err := h.images.UploadBase64(existing.ItemName, ct, *req.ImageBase64)
		if err == nil {
			existing.ImageUrl = stored
		}
	}
	updated, err := h.items.Update(id, existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to update item"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("updated", updated))
}

func (h *ItemsHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.items.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to delete item"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("deleted", gin.H{"id": id}))
}
