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
	"system-services-ui/internal/systemd"

	"github.com/gin-gonic/gin"
)

// checkUnit is invoked by the polkit rule via polkit.spawn().
// It exits 0 only if the given unit (e.g. "nginx.service") is in the config allowlist.
func checkUnit(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: system-services-ui check-unit <unit>")
		return 2
	}

	cfg, err := config.Load(config.DefaultConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		return 1
	}

	if err := systemd.SetAllowed(cfg.Services); err != nil {
		fmt.Fprintf(os.Stderr, "invalid services config: %v\n", err)
		return 1
	}

	if !systemd.IsAllowedUnit(args[0]) {
		return 1
	}
	return 0
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "check-unit" {
		os.Exit(checkUnit(os.Args[2:]))
	}

	// ---- Config ----
	cfg, err := config.Load(config.DefaultConfigPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := systemd.SetAllowed(cfg.Services); err != nil {
		log.Fatalf("invalid services config: %v", err)
	}

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
