package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IndalAwalaikal/warung-pos/backend/service"
	"github.com/gin-gonic/gin"
)

type ReportController struct{
    svc service.ReportService
}

func NewReportController(s service.ReportService) *ReportController {
    return &ReportController{svc: s}
}

// Daily returns aggregated report data for a given date (query param: date=YYYY-MM-DD). If date omitted, uses today.
func (rc *ReportController) Daily(ctx *gin.Context) {
    dateStr := ctx.Query("date")
    var d time.Time
    var err error
    if dateStr == "" {
        d = time.Now()
    } else {
        d, err = time.Parse("2006-01-02", dateStr)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid date format, use YYYY-MM-DD"})
            return
        }
    }

    out, err := rc.svc.Daily(d)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    // return the aggregated map directly to avoid fragile type assertions
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": out})
}

// Aggregate returns revenue for the last N days. Query param: days (int)
func (rc *ReportController) Aggregate(ctx *gin.Context) {
    daysStr := ctx.Query("days")
    days := 7
    if daysStr != "" {
        if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
            days = d
        }
    }
    out, err := rc.svc.Aggregate(days)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": out})
}

// ExportPDF returns a PDF report for a date (query param: date=YYYY-MM-DD)
func (rc *ReportController) ExportPDF(ctx *gin.Context) {
    dateStr := ctx.Query("date")
    var d time.Time
    var err error
    if dateStr == "" {
        d = time.Now()
    } else {
        d, err = time.Parse("2006-01-02", dateStr)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid date format"})
            return
        }
    }
    data, err := rc.svc.ExportPDF(d)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.Header("Content-Disposition", "attachment; filename=report.pdf")
    ctx.Data(http.StatusOK, "application/pdf", data)
}

// ExportExcel returns an Excel (.xlsx) report for a date
func (rc *ReportController) ExportExcel(ctx *gin.Context) {
    dateStr := ctx.Query("date")
    var d time.Time
    var err error
    if dateStr == "" {
        d = time.Now()
    } else {
        d, err = time.Parse("2006-01-02", dateStr)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid date format"})
            return
        }
    }
    data, err := rc.svc.ExportExcel(d)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.Header("Content-Disposition", "attachment; filename=report.xlsx")
    ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

