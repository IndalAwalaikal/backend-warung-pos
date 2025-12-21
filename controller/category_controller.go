package controller

import (
	"net/http"

	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/service"
	"github.com/gin-gonic/gin"
)

type CategoryController struct{
    svc service.CategoryService
}

func NewCategoryController(s service.CategoryService) *CategoryController {
    return &CategoryController{svc: s}
}

func (c *CategoryController) Create(ctx *gin.Context) {
    var in model.Category
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

func (c *CategoryController) List(ctx *gin.Context) {
    list, err := c.svc.List()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"data": list})
}
