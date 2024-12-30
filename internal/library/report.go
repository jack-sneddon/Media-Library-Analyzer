package library

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//go:embed templates/*
var templateFS embed.FS

type Report struct {
	result *AnalysisResult
}

type YearEntry struct {
	Year int
	Data *YearData
}

type TemplateData struct {
	TotalFiles  int
	LastUpdated time.Time
	Years       map[int]*YearData
	SortedYears []YearEntry
}

func NewReport(result *AnalysisResult) *Report {
	return &Report{result: result}
}

func (r *Report) GenerateSummary() string {
	var summary strings.Builder
	currentYear := time.Now().Year()

	fmt.Fprintf(&summary, "📊 Analysis Summary\n")
	fmt.Fprintf(&summary, "Total Files: %d\n", r.result.TotalFiles)
	fmt.Fprintf(&summary, "Years Covered: %d\n\n", len(r.result.Years))

	// Get issues grouped by priority
	highPriority, mediumPriority, lowPriority := r.getPrioritizedIssues(currentYear)

	// High Priority (Recent years with missing months)
	if len(highPriority) > 0 {
		fmt.Fprintf(&summary, "🚨 High Priority Issues (Last 2 years):\n")
		for _, issue := range highPriority {
			if issue.FileCount == 0 {
				fmt.Fprintf(&summary, "  • Missing: %s %d\n", issue.Month, issue.Year)
			} else {
				fmt.Fprintf(&summary, "  • Light: %s %d (%d files)\n", issue.Month, issue.Year, issue.FileCount)
			}
		}
		fmt.Fprintf(&summary, "\n")
	}

	// Medium Priority (3-5 years ago)
	if len(mediumPriority) > 0 {
		fmt.Fprintf(&summary, "⚠️  Medium Priority Issues (3-5 years ago):\n")
		for _, issue := range mediumPriority {
			if issue.FileCount == 0 {
				fmt.Fprintf(&summary, "  • Missing: %s %d\n", issue.Month, issue.Year)
			} else {
				fmt.Fprintf(&summary, "  • Light: %s %d (%d files)\n", issue.Month, issue.Year, issue.FileCount)
			}
		}
		fmt.Fprintf(&summary, "\n")
	}

	// Historical issues
	if len(lowPriority) > 0 {
		fmt.Fprintf(&summary, "📝 Historical Gaps (Older than 5 years):\n")
		currentDecade := 0
		for _, issue := range lowPriority {
			decade := (issue.Year / 10) * 10
			if decade != currentDecade {
				if currentDecade != 0 {
					fmt.Fprintf(&summary, "\n")
				}
				fmt.Fprintf(&summary, "  %ds:\n", decade)
				currentDecade = decade
			}
			if issue.FileCount == 0 {
				fmt.Fprintf(&summary, "    • Missing: %s %d\n", issue.Month, issue.Year)
			} else {
				fmt.Fprintf(&summary, "    • Light: %s %d (%d files)\n", issue.Month, issue.Year, issue.FileCount)
			}
		}
		fmt.Fprintf(&summary, "\n")
	}

	// Add recommendations
	fmt.Fprintf(&summary, "📋 Recommendations:\n")
	if len(highPriority) > 0 {
		fmt.Fprintf(&summary, "1. Focus on recent gaps first (last 2 years)\n")
	}
	if len(mediumPriority) > 0 {
		fmt.Fprintf(&summary, "2. Then address medium-term gaps (3-5 years ago)\n")
	}
	fmt.Fprintf(&summary, "3. Consider batch-processing historical gaps by decade\n")
	if r.result.TotalFiles > 0 {
		fmt.Fprintf(&summary, "4. Check 'possible_duplicates' folder for missing content\n")
	}

	return summary.String()
}

type PriorityIssue struct {
	Year      int
	Month     string
	FileCount int
}

func (r *Report) getPrioritizedIssues(currentYear int) ([]PriorityIssue, []PriorityIssue, []PriorityIssue) {
	var high, medium, low []PriorityIssue

	for year, yearData := range r.result.Years {
		for month, monthData := range yearData.Months {
			// Skip current year's future months
			if year == currentYear {
				currentMonth := time.Now().Month().String()
				if month > currentMonth {
					continue
				}
			}

			// Only process missing or light months
			if monthData.Status == "normal" {
				continue
			}

			issue := PriorityIssue{
				Year:      year,
				Month:     month,
				FileCount: monthData.FileCount,
			}

			age := currentYear - year
			switch {
			case age <= 2:
				high = append(high, issue)
			case age <= 5:
				medium = append(medium, issue)
			default:
				low = append(low, issue)
			}
		}
	}

	// Sort each priority list
	for _, list := range [][]PriorityIssue{high, medium, low} {
		sort.Slice(list, func(i, j int) bool {
			if list[i].Year != list[j].Year {
				return list[i].Year > list[j].Year // Recent years first
			}
			return list[i].Month < list[j].Month
		})
	}

	return high, medium, low
}

func (r *Report) SaveHTML(filename string) error {
	// Create sorted years slice
	var sortedYears []YearEntry
	for year, data := range r.result.Years {
		sortedYears = append(sortedYears, YearEntry{year, data})
	}
	sort.Slice(sortedYears, func(i, j int) bool {
		return sortedYears[i].Year > sortedYears[j].Year
	})

	// Prepare template data
	data := TemplateData{
		TotalFiles:  r.result.TotalFiles,
		LastUpdated: r.result.LastUpdated,
		Years:       r.result.Years,
		SortedYears: sortedYears,
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(filename)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Parse template
	tmpl, err := template.ParseFS(templateFS, "templates/report.html")
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	// Create output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	// Extract static files
	if err := r.extractStaticFiles(outputDir); err != nil {
		return fmt.Errorf("extracting static files: %w", err)
	}

	return nil
}

func (r *Report) extractStaticFiles(outputDir string) error {
	// Walk through embedded static files
	return fs.WalkDir(templateFS, "templates/static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root static directory
		if path == "templates/static" {
			return nil
		}

		// Create relative path
		relPath := path[len("templates/"):]
		fullPath := filepath.Join(outputDir, relPath)

		// Create directories
		if d.IsDir() {
			return os.MkdirAll(fullPath, 0755)
		}

		// Read and write files
		content, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			return fmt.Errorf("writing file %s: %w", fullPath, err)
		}

		return nil
	})
}
