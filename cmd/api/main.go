package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ekastn/commute-analyzer/internal/config"
	"github.com/ekastn/commute-analyzer/internal/env"
	"github.com/ekastn/commute-analyzer/internal/handler"
	"github.com/ekastn/commute-analyzer/internal/service"
	"github.com/ekastn/commute-analyzer/internal/store"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if env.GetString("SRV_ENV", "dev") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal("DB connection:", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	queries := store.New(pool)
	ors := service.NewORSClient(cfg.ORSAPIKey, httpClient)
	userService := service.NewUserService(queries)
	svc := service.NewCommuteService(queries, ors, userService)
	h := handler.NewCommuteHandler(svc)

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")
	{
		v1.POST("/commutes", h.CreateCommute)
		v1.GET("/commutes", h.ListCommutes)
		v1.PATCH("/commutes/:id", h.UpdateCommute)
		v1.DELETE("/commutes/:id", h.DeleteCommute)
	}

	srv := &http.Server{Addr: cfg.Addr, Handler: r}

	go func() {
		log.Printf("Server on %s", cfg.Addr)
		srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
