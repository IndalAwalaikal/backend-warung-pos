package router

import (
	"time"

	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/controller"
	"github.com/IndalAwalaikal/warung-pos/backend/middleware"
	crepo "github.com/IndalAwalaikal/warung-pos/backend/repository"
	cservice "github.com/IndalAwalaikal/warung-pos/backend/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    // configure CORS to allow frontend origin (use FRONTEND_ORIGIN env or default to http://localhost:5554)
    frontendOrigin := config.GetEnv("FRONTEND_ORIGIN", "http://localhost:5554")
    corsCfg := cors.Config{
        AllowOrigins:     []string{frontendOrigin},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    r.Use(cors.New(corsCfg))
    // serve uploaded files
    r.Static("/uploads", "./uploads")

    // repositories
    userRepo := crepo.NewUserRepository()
    catRepo := crepo.NewCategoryRepository()
    menuRepo := crepo.NewMenuRepository()
    txRepo := crepo.NewTransactionRepository()

    // services
    authSvc := cservice.NewAuthService(userRepo)
    catSvc := cservice.NewCategoryService(catRepo)
    menuSvc := cservice.NewMenuService(menuRepo)
    txSvc := cservice.NewTransactionService(txRepo)
    reportSvc := cservice.NewReportService(txRepo)

    // controllers
    authCtrl := controller.NewAuthController(authSvc)
    catCtrl := controller.NewCategoryController(catSvc)
    menuCtrl := controller.NewMenuController(menuSvc)
    txCtrl := controller.NewTransactionController(txSvc)
    uploadCtrl := controller.NewUploadController()

    api := r.Group("/api")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/register", authCtrl.Register)
            auth.POST("/login", authCtrl.Login)
        }

    // public
    api.GET("/categories", catCtrl.List)
    api.GET("/menus", menuCtrl.List)
    api.GET("/menus/:id", menuCtrl.Get)
    api.GET("/transactions", txCtrl.List)
    api.GET("/transactions/:id", txCtrl.Get)

        // protected: need auth
        authRequired := api.Group("")
        authRequired.Use(middleware.AuthRequired(userRepo))
        {
            authRequired.GET("/auth/me", authCtrl.Me)
            authRequired.POST("/transactions", txCtrl.Create)
                // notifications (SSE)
                notifCtrl := controller.NewNotificationController()
                authRequired.GET("/notifications/stream", notifCtrl.Stream)
            // reports
            reportCtrl := controller.NewReportController(reportSvc)
            authRequired.GET("/reports/daily", reportCtrl.Daily)
            authRequired.GET("/reports/aggregate", reportCtrl.Aggregate)
            authRequired.GET("/reports/daily/pdf", reportCtrl.ExportPDF)
            authRequired.GET("/reports/daily/excel", reportCtrl.ExportExcel)
            // toggle availability for menu (admin could also use update)
            authRequired.PATCH("/menus/:id/availability", menuCtrl.Update)
            authRequired.POST("/uploads", uploadCtrl.Upload)
        }

        // admin-only routes
        admin := api.Group("")
        admin.Use(middleware.AuthRequired(userRepo), middleware.AdminRequired())
        {
            admin.POST("/categories", catCtrl.Create)
            admin.POST("/menus", menuCtrl.Create)
            admin.PUT("/menus/:id", menuCtrl.Update)
            admin.DELETE("/menus/:id", menuCtrl.Delete)
        }
    }

    return r
}
