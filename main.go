package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

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
	//router
	router := chi.NewRouter()
	//register handler
	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", "https://*"},
			AllowedMethods: []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{"*"},
		}),
	)
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
