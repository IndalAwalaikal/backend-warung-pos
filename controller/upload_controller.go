package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadController struct{}

func NewUploadController() *UploadController { return &UploadController{} }

// Upload handles multipart file upload (field name: file)
// returns URL path for the uploaded file
func (u *UploadController) Upload(ctx *gin.Context) {
    file, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status":"error","message":"file is required"})
        return
    }

    // ensure uploads dir exists
    uploadDir := filepath.Join(".", "uploads")
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message":"cannot create upload dir"})
        return
    }

    // sanitize and create unique filename
    ext := filepath.Ext(file.Filename)
    name := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    dst := filepath.Join(uploadDir, name)

    if err := ctx.SaveUploadedFile(file, dst); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status":"error","message":"failed to save file"})
        return
    }

    // return public path
    publicURL := "/uploads/" + name
    ctx.JSON(http.StatusCreated, gin.H{"status":"success","data": gin.H{"url": publicURL}})
}
