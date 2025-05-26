package main

import (
	_ "calculator/docs"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Calculator App API
// @version 1.0
// @description API Server for Calculator Application

//@host localhost:8080
// @BasePath /

func main() {
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/compute", handleCompute)
	http.HandleFunc("/jrpc", handleJRPC)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
