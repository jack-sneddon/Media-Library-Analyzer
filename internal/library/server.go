package library

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type apiResponse struct {
	Years       map[int]*YearData `json:"years"`
	TotalFiles  int               `json:"totalFiles"`
	LastUpdated string            `json:"lastUpdated"`
}

// StartWebServer starts a local web server to serve the report
func StartWebServer(result *AnalysisResult, port int) error {
	// Create a temporary directory for the report
	tempDir, err := os.MkdirTemp("", "media-analyzer-*")
	if err != nil {
		return fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Generate the report in the temp directory
	report := NewReport(result)
	reportPath := filepath.Join(tempDir, "index.html")
	if err := report.SaveHTML(reportPath); err != nil {
		return fmt.Errorf("generating report: %w", err)
	}

	// Create file server for static files
	fileServer := http.FileServer(http.Dir(tempDir))

	// Handle static file requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})

	// Handle API requests
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Create API response
		response := apiResponse{
			Years:       result.Years,
			TotalFiles:  result.TotalFiles,
			LastUpdated: result.LastUpdated.Format("2006-01-02 15:04:05"),
		}

		// Set headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		// Encode and send response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Handle API status request
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		status := map[string]interface{}{
			"status":      "ok",
			"totalFiles":  result.TotalFiles,
			"lastUpdated": result.LastUpdated.Format("2006-01-02 15:04:05"),
			"yearsCount":  len(result.Years),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server at http://localhost%s", addr)
	log.Printf("API endpoints available at:")
	log.Printf("  - http://localhost%s/api/data", addr)
	log.Printf("  - http://localhost%s/api/status", addr)
	return http.ListenAndServe(addr, nil)
}
