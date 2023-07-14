package main

import (
	"database/sql"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	//load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("an error happen when loading the .env file: %v", err)
	}
	//now we can get the PORT attr in the current environment
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not found in the environment")
	}

	//database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL is not found in the environment")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Fail to connect to database,err: %v", err)
	}
	//get database queries
	queries := database.New(conn)
	apiCfg := apiConfig{
		DB: queries,
	}

	//router
	router := chi.NewRouter()
	//register handler
	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", "https://*"},
			AllowedMethods: []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{"*"},
			MaxAge:         300,
		}),
	)

	//v1 router
	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	router.Mount("/v1", v1Router)

	//http server
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Println("server is running on port:", port)
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatalf("an error happen when starting the http server: %v", err)
	}
}
