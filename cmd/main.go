package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project1/cmd/config"
	"project1/internal/controller"
	"project1/internal/repository"
	"project1/internal/service"
	"syscall"
	"time"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", config.User, config.Password, config.DbName, config.Host, config.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error to connect to the database: %v", err)
	}

	err = waitForDatabase(db)
	if err != nil {
		log.Fatalf("error to connect to the database: %v", err)
	}
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("error to  not run migrations: %v", err)
	}

	rep := repository.NewPostgresRep(db)

	//использую  как заглушку ,так как ручки , указанные в тз, не должны отвечать за существование юзеров
	err = rep.PopUsersIfEmpty(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	financialOperator := service.NewFinancialOperator(*rep)
	handler := controller.NewOperationController(financialOperator)
	r := gin.Default()

	r.POST("/deposit", handler.HandleDeposit)
	r.POST("/transfer", handler.HandleTransfer)
	r.GET("/transactions/:user_id", handler.HandleGetTransactions)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "/app/migrations"

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations done")
	return nil
}

func waitForDatabase(db *sql.DB) error {
	for {
		if err := db.Ping(); err == nil {
			log.Println("Database is ready")
			return nil
		}
		time.Sleep(2 * time.Second)
	}
}
