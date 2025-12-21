package controller

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/IndalAwalaikal/warung-pos/backend/dto"
	"github.com/IndalAwalaikal/warung-pos/backend/middleware"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/repository"
	"github.com/IndalAwalaikal/warung-pos/backend/service"
	"github.com/IndalAwalaikal/warung-pos/backend/utils"
	"github.com/gin-gonic/gin"
)

type TransactionController struct{
    svc service.TransactionService
}

func NewTransactionController(s service.TransactionService) *TransactionController {
    return &TransactionController{svc: s}
}

func (c *TransactionController) Create(ctx *gin.Context) {
    var req dto.TransactionCreateRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message": err.Error()})
        return
    }

    // build model.Transaction
    tx := model.Transaction{
        Subtotal:     req.Subtotal,
        Tax:          req.Tax,
        Discount:     req.Discount,
        Total:        req.Total,
        PaymentMethod: req.PaymentMethod,
        AmountPaid:   req.AmountPaid,
    }

    // attach cashier from context if available
    if v, exists := ctx.Get(middleware.ContextUserKey); exists {
        if u, ok := v.(*model.User); ok {
            tx.CashierID = &u.ID
        }
    }

    // map items
    // validate items against menu prices (prevent client price tampering)
    menuRepo := repository.NewMenuRepository()
    sum := 0.0
    for _, it := range req.Items {
        if it.Quantity <= 0 {
            ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid quantity"})
            return
        }
        m, err := menuRepo.GetByID(it.MenuID)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
            return
        }
        if m == nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message": fmt.Sprintf("menu id %d not found", it.MenuID)})
            return
        }
        // use server price
        price := m.Price
        mid := it.MenuID
        mi := model.TransactionItem{
            MenuID:   &mid,
            Quantity: it.Quantity,
            Price:    price,
        }
        tx.Items = append(tx.Items, mi)
        sum += float64(it.Quantity) * price
    }

    // calculate subtotal/tax/total if client didn't provide or to ensure integrity
    tx.Subtotal = sum
    if tx.Tax == 0 {
        tx.Tax = 0 // default tax (no tax) - change if you have tax rules
    }
    if tx.Total == 0 {
        tx.Total = tx.Subtotal + tx.Tax - tx.Discount
    }

    // default payment method
    if tx.PaymentMethod == "" {
        tx.PaymentMethod = "tunai"
    }

    if err := c.svc.Create(&tx); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"status":"success","data": tx})

    // notify connected clients about new transaction
    notif := map[string]interface{}{
        "type": "transaction_created",
        "id": tx.ID,
        "total": tx.Total,
        "cashier_id": tx.CashierID,
    }
    if b, err := json.Marshal(notif); err == nil {
        utils.NotifierInstance.Notify(string(b))
    }
}

func (c *TransactionController) List(ctx *gin.Context) {
    list, err := c.svc.List()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": list})
}

func (c *TransactionController) Get(ctx *gin.Context) {
    idStr := ctx.Param("id")
    var id uint
    if _, err := fmt.Sscan(idStr, &id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid id"})
        return
    }
    t, err := c.svc.GetByID(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    if t == nil {
        ctx.JSON(http.StatusNotFound, gin.H{"status":"error","message":"not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": t})
}
