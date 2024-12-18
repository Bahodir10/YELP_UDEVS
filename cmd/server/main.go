package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"YALP/internal/config"
	"YALP/internal/database"
	"YALP/internal/handler"
	"YALP/internal/middleware"
	"YALP/internal/repository"
	"YALP/internal/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

)

func main() {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations


	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	businessRepo := repository.NewBusinessRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize services
	userSvc := service.NewUserService(userRepo, cfg.JWTSecret)
	businessSvc := service.NewBusinessService(businessRepo, userRepo)
	reviewSvc := service.NewReviewService(reviewRepo, businessRepo)
	authSvc := service.NewAuthService(cfg.JWTSecret) // Initialize AuthService

	// Initialize handlers
	authHandler := handler.NewAuthHandler(userSvc)
	businessHandler := handler.NewBusinessHandler(businessSvc)
	reviewHandler := handler.NewReviewHandler(reviewSvc)

	// Routes
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.Logging)
	r.Use(middleware.RecoverPanic)

	// Auth routes
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")

	// Public business routes
	r.HandleFunc("/api/v1/businesses", businessHandler.ListBusinesses).Methods("GET")
	r.HandleFunc("/api/v1/businesses/search", businessHandler.SearchBusinesses).Methods("GET")
	r.HandleFunc("/api/v1/businesses/{id}", businessHandler.GetBusiness).Methods("GET")

	// Public reviews route
	r.HandleFunc("/api/v1/businesses/{id}/reviews", reviewHandler.ListReviewsForBusiness).Methods("GET")

	// Protected routes
	authenticated := r.PathPrefix("/api/v1").Subrouter()
	authenticated.Use(middleware.Auth(authSvc)) // Use AuthService for authentication middleware

	// Protected: Business
	authenticated.HandleFunc("/businesses", businessHandler.CreateBusiness).Methods("POST")
	authenticated.HandleFunc("/businesses/{id}/claim", businessHandler.ClaimBusiness).Methods("POST")

	// Protected: Reviews
	authenticated.HandleFunc("/businesses/{id}/reviews", reviewHandler.CreateReview).Methods("POST")

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// runMigrations ensures all database migrations are applied
func runMigrations(databaseURL string) error {
	log.Println("Running database migrations...")

	// Connect to the database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create migration driver: %v", err)
	}

	// Specify the migrations directory
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations", // Directory with migration files
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("unable to initialize migration: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}
