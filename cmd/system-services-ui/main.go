package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	system_services_ui "system-services-ui"
	"system-services-ui/internal/api"
	"system-services-ui/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(gin.Recovery())

	sub, err := fs.Sub(system_services_ui.StaticFS, "web/static")
	if err != nil {
		panic(err)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api.RegisterRoutes(r.Group("/api/v1"))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/ui/")
	})

	ui := r.Group("/ui")
	{
		ui.StaticFS("/", http.FS(sub))

	}

	// ---- Config ----
	cfg, err := config.Load(config.DefaultConfigPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Listen.Address, cfg.Listen.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down system-services-ui")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
