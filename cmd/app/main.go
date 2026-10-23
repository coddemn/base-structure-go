package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/coddemn/base-structure-go/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/coddemn/base-structure-go/internal/api/handler"
	"github.com/coddemn/base-structure-go/internal/config"
	"github.com/coddemn/base-structure-go/internal/database"
	"github.com/gin-gonic/gin"
)

// @title Base Golang API
// @version 1.0.0
// @description Base API-structure template on Go.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @basePath /
// @schemes http
func main() {

	//================================
	//	LOAD CONFIGS
	//================================

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	//======================
	//	CONNECT TO DB
	//======================

	ctx := context.Background()
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() {
		db.Close()
		log.Println("db connection is closed")
	}()

	log.Println("connected to database")

	//=======================
	//	INITIALYZE
	//=======================

	handler := handler.NewHandler()

	httpAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	//==============================
	//	MIDDLEWARES
	//==============================

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	//==============================
	//	ENDPOINTS
	//==============================

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hi from base-structure, Demidos!",
		})
	})
	router.GET("/health", handler.Health)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//===============================
	//	START SERVER
	//===============================

	go func() {
		log.Printf("HTTP Server is running on %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server error: %v", err)
		}
	}()

	//=================================
	//	GRACEFUL SHUTDOWN (CTRL+C)
	//=================================

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("Given signal %v. Stopping work...\n", sig)

	shtdCtx, shtdCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shtdCancel()

	if err := httpServer.Shutdown(shtdCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}

	// wait workers is finalized
	//wg.Wait()
	log.Println("All workers is completed. App correct stopped.")

}
