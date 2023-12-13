package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/joho/godotenv"
	"github.com/rigved-desai/paryatan-backend/api"
	"github.com/rigved-desai/paryatan-backend/db"
)

func main() {

	godotenv.Load()

	dbpool, err:= db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.DB.Close()
	
	
	log.Println("Connected to DB!")
	
	router := api.NewRouter(dbpool)
	c := cors.Default()
	handler := c.Handler(router)
	
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable not found.")
	}

	srv := &http.Server{
		Handler: handler,
		Addr: ":" + portString,
	}
	log.Printf("Server is running on PORT: %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}