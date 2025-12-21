package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/IndalAwalaikal/warung-pos/backend/repository"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportService interface {
    Daily(date time.Time) (map[string]interface{}, error)
    Aggregate(days int) ([]map[string]interface{}, error)
    ExportExcel(date time.Time) ([]byte, error)
    ExportPDF(date time.Time) ([]byte, error)
}

type reportService struct{
    txRepo repository.TransactionRepository
}

func NewReportService(tx repository.TransactionRepository) ReportService {
    return &reportService{txRepo: tx}
}

// Daily returns aggregated report data for the given date
func (s *reportService) Daily(date time.Time) (map[string]interface{}, error) {
    list, err := s.txRepo.List()
    if err != nil {
        return nil, err
    }

    // filter by date
    start := time.Date(date.Year(), date.Month(), date.Day(), 0,0,0,0, date.Location())
    end := start.Add(24*time.Hour)

    var totalRevenue float64
    var totalTransactions int
    var totalItems int
    menuCount := map[uint]*struct{ Name string; Count int; Revenue float64 }{}

    for _, t := range list {
        if t.CreatedAt.Before(start) || !t.CreatedAt.Before(end) {
            continue
        }
        totalTransactions++
        totalRevenue += t.Total
        for _, it := range t.Items {
            totalItems += it.Quantity
            mid := uint(0)
            if it.MenuID != nil {
                mid = *it.MenuID
            }
            if _, ok := menuCount[mid]; !ok {
                menuCount[mid] = &struct{ Name string; Count int; Revenue float64 }{Name: it.Menu.Name}
            }
            menuCount[mid].Count += it.Quantity
            menuCount[mid].Revenue += float64(it.Quantity) * it.Price
        }
    }

    best := []map[string]interface{}{}
    for k, v := range menuCount {
        best = append(best, map[string]interface{}{"id": k, "name": v.Name, "count": v.Count, "revenue": v.Revenue})
    }

    return map[string]interface{}{
        "date": start.Format("2006-01-02"),
        "total_revenue": totalRevenue,
        "total_transactions": totalTransactions,
        "total_items": totalItems,
        "best_sellers": best,
    }, nil
}

// Aggregate returns revenue per day for the last `days` days (including today)
func (s *reportService) Aggregate(days int) ([]map[string]interface{}, error) {
    list, err := s.txRepo.List()
    if err != nil {
        return nil, err
    }

    res := make([]map[string]interface{}, 0, days)
    today := time.Now()
    // for each day from days-1 .. 0
    for i := days - 1; i >= 0; i-- {
        d := time.Date(today.Year(), today.Month(), today.Day()-i, 0, 0, 0, 0, today.Location())
        start := d
        end := start.Add(24 * time.Hour)
        var revenue float64
        var txCount int
        for _, t := range list {
            if t.CreatedAt.Before(start) || !t.CreatedAt.Before(end) {
                continue
            }
            revenue += t.Total
            txCount++
        }
        res = append(res, map[string]interface{}{
            "date": start.Format("2006-01-02"),
            "revenue": revenue,
            "transactions": txCount,
        })
    }
    return res, nil
}

// ExportExcel generates an Excel file for the given date (daily report)
func (s *reportService) ExportExcel(date time.Time) ([]byte, error) {
    // reuse Daily aggregation
    daily, err := s.Daily(date)
    if err != nil {
        return nil, err
    }

    f := excelize.NewFile()
    // Summary sheet
    sheet := "Summary"
    f.NewSheet(sheet)
    f.SetCellValue(sheet, "A1", "Date")
    f.SetCellValue(sheet, "B1", daily["date"])
    f.SetCellValue(sheet, "A2", "Total Revenue")
    f.SetCellValue(sheet, "B2", daily["total_revenue"])
    f.SetCellValue(sheet, "A3", "Total Transactions")
    f.SetCellValue(sheet, "B3", daily["total_transactions"])
    f.SetCellValue(sheet, "A4", "Total Items")
    f.SetCellValue(sheet, "B4", daily["total_items"])

    // Best sellers sheet
    bsSheet := "Best Sellers"
    f.NewSheet(bsSheet)
    f.SetCellValue(bsSheet, "A1", "ID")
    f.SetCellValue(bsSheet, "B1", "Name")
    f.SetCellValue(bsSheet, "C1", "Count")
    f.SetCellValue(bsSheet, "D1", "Revenue")
    best, _ := daily["best_sellers"].([]map[string]interface{})
    // fallback: try casting from []interface{}
    if best == nil {
        best = []map[string]interface{}{}
        if raw, ok := daily["best_sellers"].([]interface{}); ok {
            for _, r := range raw {
                if m, ok := r.(map[string]interface{}); ok {
                    best = append(best, m)
                }
            }
        }
    }
    row := 2
    for _, b := range best {
        f.SetCellValue(bsSheet, fmt.Sprintf("A%d", row), b["id"])
        f.SetCellValue(bsSheet, fmt.Sprintf("B%d", row), b["name"])
        f.SetCellValue(bsSheet, fmt.Sprintf("C%d", row), b["count"])
        f.SetCellValue(bsSheet, fmt.Sprintf("D%d", row), b["revenue"])
        row++
    }

    // Set active sheet to Summary
    if idx, err := f.GetSheetIndex(sheet); err == nil {
        f.SetActiveSheet(idx)
    }

    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

// ExportPDF generates a professional PDF report with proper header, footer, and formatting
func (s *reportService) ExportPDF(date time.Time) ([]byte, error) {
    daily, err := s.Daily(date)
    if err != nil {
        return nil, err
    }

    pdf := gofpdf.New("P", "mm", "A4", "")
    
    // Add header and footer callback
    pdf.SetHeaderFunc(func() {
        // Top border line
        pdf.SetDrawColor(0, 0, 0)
        pdf.SetLineWidth(0.5)
        pdf.Line(10, 10, 200, 10)
        
        // Header title
        pdf.SetY(12)
        pdf.SetFont("Helvetica", "B", 18)
        pdf.SetTextColor(0, 0, 0)
        pdf.CellFormat(0, 8, "LAPORAN SISTEM KASIR", "", 0, "C", false, 0, "")
        pdf.Ln(6)
        pdf.SetFont("Helvetica", "B", 16)
        pdf.CellFormat(0, 8, "WARUNG MAKAN", "", 0, "C", false, 0, "")
        pdf.Ln(8)
        
        // Bottom border line after header
        pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
        pdf.Ln(3)
    })
    
    pdf.SetFooterFunc(func() {
        pdf.SetY(-15)
        pdf.SetDrawColor(0, 0, 0)
        pdf.SetLineWidth(0.3)
        pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
        pdf.Ln(2)
        pdf.SetFont("Helvetica", "I", 8)
        pdf.SetTextColor(100, 100, 100)
        pdf.CellFormat(0, 5, fmt.Sprintf("Halaman %d", pdf.PageNo()), "", 0, "C", false, 0, "")
        pdf.Ln(2)
        pdf.CellFormat(0, 5, fmt.Sprintf("Dicetak pada: %s", time.Now().Format("02 January 2006 15:04:05")), "", 0, "C", false, 0, "")
    })
    
    pdf.AddPage()
    pdf.SetTextColor(0, 0, 0)
    
    // Document information box
    pdf.SetFont("Helvetica", "B", 12)
    pdf.SetFillColor(240, 240, 240)
    pdf.SetDrawColor(200, 200, 200)
    pdf.CellFormat(0, 8, "INFORMASI LAPORAN", "1", 0, "L", true, 0, "")
    pdf.Ln(8)
    
    // Date and period
    pdf.SetFont("Helvetica", "", 11)
    pdf.SetFillColor(250, 250, 250)
    pdf.CellFormat(60, 7, "Tanggal Laporan", "1", 0, "L", true, 0, "")
    pdf.CellFormat(130, 7, fmt.Sprintf(": %s", daily["date"]), "1", 0, "L", false, 0, "")
    pdf.Ln(7)
    pdf.CellFormat(60, 7, "Periode", "1", 0, "L", true, 0, "")
    pdf.CellFormat(130, 7, ": Harian (Daily Report)", "1", 0, "L", false, 0, "")
    pdf.Ln(12)
    
    // Summary section
    pdf.SetFont("Helvetica", "B", 12)
    pdf.SetFillColor(70, 130, 180)
    pdf.SetTextColor(255, 255, 255)
    pdf.CellFormat(0, 9, "RINGKASAN PENJUALAN", "1", 0, "C", true, 0, "")
    pdf.Ln(9)
    pdf.SetTextColor(0, 0, 0)
    
    // Summary data in a box
    pdf.SetFont("Helvetica", "", 11)
    pdf.SetFillColor(245, 245, 245)
    
    // Row 1
    pdf.CellFormat(95, 8, "Total Pendapatan (Rp)", "1", 0, "L", true, 0, "")
    pdf.SetFont("Helvetica", "B", 11)
    pdf.CellFormat(95, 8, fmt.Sprintf("Rp %.2f", daily["total_revenue"]), "1", 0, "R", false, 0, "")
    pdf.Ln(8)
    
    // Row 2
    pdf.SetFont("Helvetica", "", 11)
    pdf.CellFormat(95, 8, "Jumlah Transaksi", "1", 0, "L", true, 0, "")
    pdf.SetFont("Helvetica", "B", 11)
    pdf.CellFormat(95, 8, fmt.Sprintf("%v transaksi", daily["total_transactions"]), "1", 0, "R", false, 0, "")
    pdf.Ln(8)
    
    // Row 3
    pdf.SetFont("Helvetica", "", 11)
    pdf.CellFormat(95, 8, "Total Item Terjual", "1", 0, "L", true, 0, "")
    pdf.SetFont("Helvetica", "B", 11)
    pdf.CellFormat(95, 8, fmt.Sprintf("%v item", daily["total_items"]), "1", 0, "R", false, 0, "")
    pdf.Ln(8)
    
    // Average per transaction
    avgPerTx := 0.0
    if txCount, ok := daily["total_transactions"].(int); ok && txCount > 0 {
        if rev, ok := daily["total_revenue"].(float64); ok {
            avgPerTx = rev / float64(txCount)
        }
    }
    pdf.SetFont("Helvetica", "", 11)
    pdf.CellFormat(95, 8, "Rata-rata per Transaksi (Rp)", "1", 0, "L", true, 0, "")
    pdf.SetFont("Helvetica", "B", 11)
    pdf.CellFormat(95, 8, fmt.Sprintf("Rp %.2f", avgPerTx), "1", 0, "R", false, 0, "")
    pdf.Ln(14)
    
    // Best sellers section
    pdf.SetFont("Helvetica", "B", 12)
    pdf.SetFillColor(70, 130, 180)
    pdf.SetTextColor(255, 255, 255)
    pdf.CellFormat(0, 9, "DAFTAR MENU TERLARIS", "1", 0, "C", true, 0, "")
    pdf.Ln(9)
    pdf.SetTextColor(0, 0, 0)
    
    // Table header
    pdf.SetFont("Helvetica", "B", 10)
    pdf.SetFillColor(220, 220, 220)
    pdf.CellFormat(15, 8, "No", "1", 0, "C", true, 0, "")
    pdf.CellFormat(90, 8, "Nama Menu", "1", 0, "L", true, 0, "")
    pdf.CellFormat(30, 8, "Jumlah", "1", 0, "C", true, 0, "")
    pdf.CellFormat(55, 8, "Pendapatan (Rp)", "1", 0, "R", true, 0, "")
    pdf.Ln(8)
    
    // Parse best sellers
    var best []map[string]interface{}
    if bs, ok := daily["best_sellers"].([]map[string]interface{}); ok {
        best = bs
    } else if raw, ok := daily["best_sellers"].([]interface{}); ok {
        for _, r := range raw {
            if m, ok := r.(map[string]interface{}); ok {
                best = append(best, m)
            }
        }
    }
    
    // Table rows
    pdf.SetFont("Helvetica", "", 10)
    for i, b := range best {
        name, _ := b["name"].(string)
        count := int64(0)
        if v, ok := b["count"].(int); ok {
            count = int64(v)
        } else if v, ok := b["count"].(float64); ok {
            count = int64(v)
        }
        revenue := 0.0
        if v, ok := b["revenue"].(float64); ok {
            revenue = v
        }
        
        // Alternating row colors
        if i%2 == 0 {
            pdf.SetFillColor(255, 255, 255)
        } else {
            pdf.SetFillColor(245, 245, 245)
        }
        
        pdf.CellFormat(15, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", true, 0, "")
        pdf.CellFormat(90, 7, name, "1", 0, "L", true, 0, "")
        pdf.CellFormat(30, 7, fmt.Sprintf("%d", count), "1", 0, "C", true, 0, "")
        pdf.CellFormat(55, 7, fmt.Sprintf("%.2f", revenue), "1", 0, "R", true, 0, "")
        pdf.Ln(7)
    }
    
    // Visual chart section
    if len(best) > 0 && len(best) <= 10 {
        pdf.Ln(8)
        pdf.SetFont("Helvetica", "B", 12)
        pdf.SetFillColor(70, 130, 180)
        pdf.SetTextColor(255, 255, 255)
        pdf.CellFormat(0, 9, "GRAFIK PENDAPATAN MENU", "1", 0, "C", true, 0, "")
        pdf.Ln(12)
        pdf.SetTextColor(0, 0, 0)
        
        // Compute max for scaling
        max := 0.0
        for _, b := range best {
            if v, ok := b["revenue"].(float64); ok && v > max {
                max = v
            }
        }
        
        if max > 0 {
            pdf.SetFont("Helvetica", "", 9)
            for _, b := range best {
                name, _ := b["name"].(string)
                revenue := 0.0
                if v, ok := b["revenue"].(float64); ok {
                    revenue = v
                }
                
                // Truncate long names
                if len(name) > 25 {
                    name = name[:22] + "..."
                }
                
                width := (revenue / max) * 120 // scale to 120mm max
                
                pdf.SetFont("Helvetica", "", 8)
                pdf.CellFormat(60, 6, name, "", 0, "L", false, 0, "")
                
                x := pdf.GetX()
                y := pdf.GetY()
                
                // Draw bar
                pdf.SetFillColor(70, 130, 180)
                pdf.Rect(x, y, width, 5, "F")
                
                // Draw border
                pdf.SetDrawColor(50, 100, 150)
                pdf.Rect(x, y, width, 5, "D")
                
                // Value label
                pdf.SetX(x + width + 2)
                pdf.SetFont("Helvetica", "B", 8)
                pdf.CellFormat(30, 6, fmt.Sprintf("Rp %.0f", revenue), "", 0, "L", false, 0, "")
                pdf.Ln(7)
            }
        }
    }
    
    // Signature section
    pdf.Ln(15)
    pdf.SetFont("Helvetica", "", 10)
    pdf.CellFormat(95, 6, "", "", 0, "C", false, 0, "")
    pdf.CellFormat(95, 6, fmt.Sprintf("%s, %s", "Makassar", time.Now().Format("02 January 2006")), "", 0, "C", false, 0, "")
    pdf.Ln(6)
    pdf.CellFormat(95, 6, "", "", 0, "C", false, 0, "")
    pdf.CellFormat(95, 6, "Manajer Warung", "", 0, "C", false, 0, "")
    pdf.Ln(18)
    pdf.CellFormat(95, 6, "", "", 0, "C", false, 0, "")
    pdf.SetFont("Helvetica", "B", 10)
    pdf.CellFormat(95, 6, "(_________________)", "", 0, "C", false, 0, "")
    
    var buf bytes.Buffer
    if err := pdf.Output(&buf); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}