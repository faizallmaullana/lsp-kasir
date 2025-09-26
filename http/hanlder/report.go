package handler

import (
	"net/http"
	"strconv"
	"time"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/services"
	"faizalmaulana/lsp/models/repo"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	cfg       *conf.Config
	txSvc     services.TransactionsService
	pivotRepo repo.PivotItemsToTransactionsRepo
}

func NewReportHandler(cfg *conf.Config, tx services.TransactionsService, pivot repo.PivotItemsToTransactionsRepo) *ReportHandler {
	return &ReportHandler{cfg: cfg, txSvc: tx, pivotRepo: pivot}
}

func (h *ReportHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/reports/date/:dd/:mm/:yyyy", h.reportByExactDate)
	rg.GET("/reports/:bulan/:tahun", h.reportByMonthYear)
	rg.GET("/reports/today", h.reportToday)
	rg.GET("/reports/today/summary", h.reportTodaySummary)
}

func (h *ReportHandler) reportByExactDate(c *gin.Context) {
	ddStr := c.Param("dd")
	mmStr := c.Param("mm")
	yyStr := c.Param("yyyy")

	day, errD := strconv.Atoi(ddStr)
	month, errM := strconv.Atoi(mmStr)
	year, errY := strconv.Atoi(yyStr)
	if errD != nil || errM != nil || errY != nil || day < 1 || day > 31 || month < 1 || month > 12 || year < 1970 {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse("invalid dd/mm/yyyy"))
		return
	}

	list, err := h.txSvc.GetAll(1000, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to query transactions"))
		return
	}

	var items []dto.ReportTransaction
	total := 0
	var sum float64

	for _, t := range list {
		ts := t.Timestamp
		if ts.IsZero() {
			continue
		}
		if ts.Year() == year && int(ts.Month()) == month && ts.Day() == day {
			items = append(items, dto.ReportTransaction{
				IdTransaction: t.IdTransaction,
				TotalPrice:    t.TotalPrice,
				BuyerContact:  t.BuyerContact,
				Timestamp:     ts.Format(time.RFC3339),
			})
			sum += t.TotalPrice
			total++
		}
	}

	resp := dto.TodayReportResponse{Date: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Format("2006-01-02"), Total: total, Sum: sum, Items: items}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", resp))
}

func (h *ReportHandler) reportByMonthYear(c *gin.Context) {
	monthStr := c.Param("bulan")
	yearStr := c.Param("tahun")
	month, err1 := strconv.Atoi(monthStr)
	year, err2 := strconv.Atoi(yearStr)
	if err1 != nil || err2 != nil || month < 1 || month > 12 || year < 1970 {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse("invalid month/year"))
		return
	}

	list, err := h.txSvc.GetAll(1000, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to query transactions"))
		return
	}

	var items []dto.ReportTransaction
	total := 0
	var sum float64
	for _, t := range list {
		ts := t.Timestamp
		if ts.IsZero() {
			continue
		}
		if int(ts.Month()) == month && ts.Year() == year {
			items = append(items, dto.ReportTransaction{
				IdTransaction: t.IdTransaction,
				TotalPrice:    t.TotalPrice,
				BuyerContact:  t.BuyerContact,
				Timestamp:     ts.Format(time.RFC3339),
			})
			sum += t.TotalPrice
			total++
		}
	}

	resp := dto.ReportResponse{Month: month, Year: year, Total: total, Sum: sum, Items: items}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", resp))
}

func (h *ReportHandler) reportToday(c *gin.Context) {
	now := time.Now()
	y, m, d := now.Date()

	list, err := h.txSvc.GetAll(1000, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to query transactions"))
		return
	}

	var items []dto.ReportTransaction
	total := 0
	var sum float64
	for _, t := range list {
		ts := t.Timestamp
		if ts.IsZero() {
			continue
		}
		ty, tm, td := ts.Date()
		if ty == y && tm == m && td == d {
			items = append(items, dto.ReportTransaction{
				IdTransaction: t.IdTransaction,
				TotalPrice:    t.TotalPrice,
				BuyerContact:  t.BuyerContact,
				Timestamp:     ts.Format(time.RFC3339),
			})
			sum += t.TotalPrice
			total++
		}
	}

	resp := dto.TodayReportResponse{Date: now.Format("2006-01-02"), Total: total, Sum: sum, Items: items}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", resp))
}

// reportTodaySummary returns today's revenue (sum), total transactions, and total products sold.
func (h *ReportHandler) reportTodaySummary(c *gin.Context) {
	now := time.Now()
	y, m, d := now.Date()

	list, err := h.txSvc.GetAll(1000, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to query transactions"))
		return
	}

	totalTx := 0
	var sum float64
	totalProducts := 0

	for _, t := range list {
		ts := t.Timestamp
		if ts.IsZero() {
			continue
		}
		ty, tm, td := ts.Date()
		if ty == y && tm == m && td == d {
			totalTx++
			sum += t.TotalPrice
			pivots, _ := h.pivotRepo.ListByTransaction(t.IdTransaction)
			for _, p := range pivots {
				totalProducts += p.Quantity
			}
		}
	}

	resp := dto.TodaySummaryResponse{
		Date:              now.Format("2006-01-02"),
		TotalTransactions: totalTx,
		TotalProductsSold: totalProducts,
		SumTotalPrice:     sum,
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", resp))
}
