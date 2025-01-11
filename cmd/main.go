package main

import (
	"log"
	"mcq-test-system/internal/handlers"
	"net/http"
	//"path/filepath"
)

func main() {
	// Initialize handler
	testHandler := handlers.NewTestHandler() // Load all HTML templates when handler is created

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler.StartTest)
	mux.HandleFunc("/question", testHandler.HandleQuestion)
	mux.HandleFunc("/submit", testHandler.HandleSubmit)
	mux.HandleFunc("/report", testHandler.GenerateReport)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
