package controller

import (
	"fmt"
	"net/http"

	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/service"
	"github.com/gin-gonic/gin"
)

type MenuController struct{
    svc service.MenuService
}

func NewMenuController(s service.MenuService) *MenuController {
    return &MenuController{svc: s}
}

func (c *MenuController) Create(ctx *gin.Context) {
    var in model.Menu
    if err := ctx.ShouldBindJSON(&in); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message": err.Error()})
        return
    }
    if err := c.svc.Create(&in); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"status":"success","data": in})
}

func (c *MenuController) List(ctx *gin.Context) {
    list, err := c.svc.List()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": list})
}

func (c *MenuController) Get(ctx *gin.Context) {
    idStr := ctx.Param("id")
    var id uint
    if _, err := fmt.Sscan(idStr, &id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid id"})
        return
    }
    m, err := c.svc.GetByID(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    if m == nil {
        ctx.JSON(http.StatusNotFound, gin.H{"status":"error","message":"not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": m})
}

func (c *MenuController) Update(ctx *gin.Context) {
    idStr := ctx.Param("id")
    var id uint
    if _, err := fmt.Sscan(idStr, &id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid id"})
        return
    }
    // fetch existing
    existing, err := c.svc.GetByID(id)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    if existing == nil {
        ctx.JSON(http.StatusNotFound, gin.H{"status":"error","message": "not found"})
        return
    }

    // bind into map to allow partial updates
    var payload map[string]interface{}
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message": err.Error()})
        return
    }

    // merge allowed fields
    if v, ok := payload["name"].(string); ok {
        existing.Name = v
    }
    if v, ok := payload["description"].(string); ok {
        existing.Description = v
    }
    if v, ok := payload["image_url"].(string); ok {
        existing.ImageURL = v
    }
    if v, ok := payload["price"].(float64); ok {
        existing.Price = v
    }
    if v, ok := payload["category_id"].(float64); ok {
        u := uint(v)
        existing.CategoryID = &u
    }
    if v, ok := payload["is_available"].(bool); ok {
        existing.IsAvailable = v
    }

    if err := c.svc.Update(existing); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": existing})
}

func (c *MenuController) Delete(ctx *gin.Context) {
    idStr := ctx.Param("id")
    var id uint
    if _, err := fmt.Sscan(idStr, &id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"invalid id"})
        return
    }
    if err := c.svc.Delete(id); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","message":"deleted"})
}
