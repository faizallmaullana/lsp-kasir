package handler

import (
	"net/http"
	"strconv"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"
	"faizalmaulana/lsp/models/entity"
	"faizalmaulana/lsp/models/repo"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type TransactionsHandler struct {
	cfg       *conf.Config
	txSvc     services.TransactionsService
	itemsRepo repo.ItemsRepo
	pivotRepo repo.PivotItemsToTransactionsRepo
}

func NewTransactionsHandler(cfg *conf.Config, tx services.TransactionsService, items repo.ItemsRepo, pivot repo.PivotItemsToTransactionsRepo) *TransactionsHandler {
	return &TransactionsHandler{cfg: cfg, txSvc: tx, itemsRepo: items, pivotRepo: pivot}
}

func (h *TransactionsHandler) Register(rr *gin.RouterGroup) {
	rg := rr.Group("/transactions")
	rg.GET("", h.list)
	rg.GET(":id", h.get)
	rg.POST("", middleware.JWTMiddleware(h.cfg), h.create)
	rg.PUT(":id", middleware.JWTMiddleware(h.cfg), h.update)
	rg.DELETE(":id", middleware.JWTMiddleware(h.cfg), h.delete)
}

func (h *TransactionsHandler) list(c *gin.Context) {
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
	out, err := h.txSvc.GetAll(count, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to list transactions"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", out))
}

func (h *TransactionsHandler) get(c *gin.Context) {
	id := c.Param("id")
	t, err := h.txSvc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("transaction not found"))
		return
	}
	// load items (optional)
	pivots, _ := h.pivotRepo.ListByTransaction(id)
	// enrich each pivot with item details
	details := make([]dto.TransactionItemDetail, 0, len(pivots))
	for _, p := range pivots {
		it, err := h.itemsRepo.GetByID(p.IdItem)
		if err != nil {
			// skip missing items, but continue
			continue
		}
		details = append(details, dto.TransactionItemDetail{
			IdItem:   it.IdItem,
			ItemName: it.ItemName,
			ImageUrl: it.ImageUrl,
			Quantity: p.Quantity,
			Price:    p.Price,
		})
	}
	resp := gin.H{"transaction": t, "items": details}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", resp))
}

func (h *TransactionsHandler) create(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse("items required"))
		return
	}

	// Extract user ID from JWT claims (set by middleware)
	userID := ""
	if v, ok := c.Get("claims"); ok {
		switch claims := v.(type) {
		case jwt.MapClaims:
			if sub, ok := claims["sub"].(string); ok {
				userID = sub
			}
		case map[string]any:
			if sub, ok := claims["sub"].(string); ok {
				userID = sub
			}
		}
	}
	if userID == "" {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}

	// Validate items and compute total
	var total float64
	pivots := make([]entity.PivotItemsToTransaction, 0, len(req.Items))
	for _, it := range req.Items {
		item, err := h.itemsRepo.GetByID(it.IdItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.BadRequestResponse("invalid item: "+it.IdItem))
			return
		}
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		pivots = append(pivots, entity.PivotItemsToTransaction{
			IdItem:   item.IdItem,
			Quantity: qty,
			Price:    item.Price,
		})
		total += float64(qty) * item.Price
	}

	tx := &entity.Transactions{
		IdTransaction: helper.Uuid(),
		IdUser:        userID,
		BuyerContact:  req.BuyerContact,
		TotalPrice:    total,
	}

	// Create transaction
	saved, err := h.txSvc.Create(tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to create transaction"))
		return
	}

	// attach transaction id and persist pivot rows
	for i := range pivots {
		pivots[i].IdTransaction = saved.IdTransaction
	}
	if err := h.pivotRepo.BulkCreate(pivots); err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to save items"))
		return
	}

	c.JSON(http.StatusCreated, helper.SuccessResponse("created", gin.H{"transaction": saved, "items": pivots}))
}

func (h *TransactionsHandler) update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	t, err := h.txSvc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("transaction not found"))
		return
	}
	if req.BuyerContact != nil {
		t.BuyerContact = *req.BuyerContact
	}
	updated, err := h.txSvc.Update(id, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to update transaction"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("updated", updated))
}

func (h *TransactionsHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.txSvc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to delete transaction"))
		return
	}
	// soft delete pivot rows too
	_ = h.pivotRepo.DeleteByTransaction(id)
	c.JSON(http.StatusOK, helper.SuccessResponse("deleted", gin.H{"id": id}))
}
