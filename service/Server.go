package service

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"shipping/routes"
	"time"
)

type Server struct {
}

func (p *Server) RunServer() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	routes := new(routes.Routes)
	envport := os.Getenv("PORT")
	port := ":" + envport
	router := mux.NewRouter()
	routes.SetupRoutes(router)
	server := &http.Server{
		Addr:              port,
		Handler:           router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
	server.ListenAndServe()
}
