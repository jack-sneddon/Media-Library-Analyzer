package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/jack-sneddon/Media-Library-Analyzer/internal/library"
)

func main() {
	// Parse command line flags
	mediaPath := flag.String("path", ".", "Path to media library")
	webMode := flag.Bool("web", false, "Start web dashboard")
	port := flag.Int("port", 8080, "Port for web dashboard")
	flag.Parse()

	// Get absolute path
	absPath, err := filepath.Abs(*mediaPath)
	if err != nil {
		log.Fatalf("Error resolving path: %v", err)
	}

	// Create analyzer instance
	analyzer := library.NewAnalyzer(absPath)

	// Run analysis
	fmt.Println("📸 Media Library Analyzer")
	fmt.Println("========================")
	fmt.Printf("Analyzing library at: %s\n\n", absPath)

	result, err := analyzer.Analyze()
	if err != nil {
		log.Fatalf("Analysis failed: %v", err)
	}

	// Generate report
	report := library.NewReport(result)
	summary := report.GenerateSummary()
	fmt.Println(summary)

	// Save detailed report
	reportPath := "library_analysis.html"
	if err := report.SaveHTML(reportPath); err != nil {
		log.Printf("Warning: Could not save HTML report: %v", err)
	} else {
		fmt.Printf("\nDetailed report saved to: %s\n", reportPath)
	}

	// Start web server if requested
	if *webMode {
		fmt.Printf("\nStarting web dashboard at http://localhost:%d\n", *port)
		fmt.Println("Press Ctrl+C to exit")
		if err := library.StartWebServer(result, *port); err != nil {
			log.Fatalf("Web server error: %v", err)
		}
	}
}
