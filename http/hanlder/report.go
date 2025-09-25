package handler

import (
	"net/http"
	"strconv"
	"time"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/dto"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	cfg   *conf.Config
	txSvc services.TransactionsService
}

func NewReportHandler(cfg *conf.Config, tx services.TransactionsService) *ReportHandler {
	return &ReportHandler{cfg: cfg, txSvc: tx}
}

func (h *ReportHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/report/date/:dd/:mm/:yyyy", h.reportByExactDate)
	rg.GET("/report/:bulan/:tahun", h.reportByMonthYear)
	rg.GET("/report/today", h.reportToday)
}

// reportByExactDate returns report filtered by an exact calendar date (dd/mm/yyyy).
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

	// Pull transactions (naive) and filter by date components
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

	// naive approach: list all and filter by month/year
	// For large data, add repo method to query by date range.
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
