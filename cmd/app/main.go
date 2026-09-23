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

	"github.com/coddemn/base-structure-go/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {

	//================================
	//	LOAD CONFIGS
	//================================

	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	//=======================
	//	INITIALYZE
	//=======================

	httpAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: router,
	}

	//==============================
	//	ENDPOINTS
	//==============================

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hi from get-access, Demidos!",
		})
	})

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
