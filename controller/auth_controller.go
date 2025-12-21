package controller

import (
	"net/http"

	"github.com/IndalAwalaikal/warung-pos/backend/middleware"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct{
    svc service.AuthService
}

func NewAuthController(s service.AuthService) *AuthController {
    return &AuthController{svc: s}
}

func (c *AuthController) Register(ctx *gin.Context) {
    var in model.User
    if err := ctx.ShouldBindJSON(&in); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := c.svc.Register(&in); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (c *AuthController) Login(ctx *gin.Context) {
    var req struct{
        Email string `json:"email"`
        Password string `json:"password"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message": err.Error()})
        return
    }
    token, user, err := c.svc.Login(req.Email, req.Password)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message": err.Error()})
        return
    }
    if token == "" || user == nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status":"error","message": "invalid credentials"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": gin.H{"token": token, "user": user}})
}

func (c *AuthController) Me(ctx *gin.Context) {
    v, exists := ctx.Get(middleware.ContextUserKey)
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status":"error","message":"unauthorized"})
        return
    }
    user, ok := v.(*model.User)
    if !ok || user == nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status":"error","message":"unauthorized"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"status":"success","data": user})
}
