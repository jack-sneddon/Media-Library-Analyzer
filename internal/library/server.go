package library

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// StartWebServer starts a local web server to display the analysis results
func StartWebServer(result *AnalysisResult, port int) error {
	// Serve static files from web directory
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// API endpoint for JSON data
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}