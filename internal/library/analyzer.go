package library

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MonthData struct {
	Name      string `json:"name"`
	FileCount int    `json:"fileCount"`
	Status    string `json:"status"` // "missing", "light", or "normal"
}

type YearData struct {
	Year   int                   `json:"year"`
	Months map[string]*MonthData `json:"months"`
}

type AnalysisResult struct {
	Years       map[int]*YearData `json:"years"`
	TotalFiles  int               `json:"totalFiles"`
	LastUpdated time.Time         `json:"lastUpdated"`
}

type Analyzer struct {
	rootPath string
}

func NewAnalyzer(rootPath string) *Analyzer {
	return &Analyzer{
		rootPath: rootPath,
	}
}

func (y *YearData) TotalFiles() int {
	total := 0
	for _, month := range y.Months {
		total += month.FileCount
	}
	return total
}

// SortedMonths returns a slice of MonthData sorted by chronological order
func (y *YearData) SortedMonths() []*MonthData {
	var months []*MonthData
	for name, data := range y.Months {
		data.Name = name
		months = append(months, data)
	}

	sort.Slice(months, func(i, j int) bool {
		return monthOrder(months[i].Name) < monthOrder(months[j].Name)
	})

	return months
}

func (a *Analyzer) extractYear(dirName string) (int, error) {
	// Handle different year patterns
	patterns := []string{
		`^(\d{4})-.*`,      // e.g., "1987- Jack"
		`^(\d{4})$`,        // e.g., "2023"
		`^(\d{4}) - .*`,    // e.g., "1994 - Lynsey"
		`^(\d{4})-\d{4}.*`, // e.g., "1989-1990 - Jack"
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(dirName); len(matches) > 1 {
			return strconv.Atoi(matches[1])
		}
	}

	return 0, fmt.Errorf("no valid year found in directory name: %s", dirName)
}

func (a *Analyzer) Analyze() (*AnalysisResult, error) {
	result := &AnalysisResult{
		Years:       make(map[int]*YearData),
		LastUpdated: time.Now(),
	}

	// Read root directory
	entries, err := ioutil.ReadDir(a.rootPath)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %v", a.rootPath, err)
	}

	// Process each year directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip special directories
		if entry.Name() == "possible_duplicates" ||
			!strings.Contains(entry.Name(), "-") && len(entry.Name()) != 4 {
			continue
		}

		// Parse year from directory name
		year, err := a.extractYear(entry.Name())
		if err != nil {
			continue // Skip invalid year directories
		}

		yearData := &YearData{
			Year:   year,
			Months: make(map[string]*MonthData),
		}

		// Initialize all months
		for i := 1; i <= 12; i++ {
			monthName := time.Month(i).String()
			yearData.Months[monthName] = &MonthData{
				Name:      monthName,
				FileCount: 0,
				Status:    "missing",
			}
		}

		// Process month directories
		yearPath := filepath.Join(a.rootPath, entry.Name())
		monthEntries, err := ioutil.ReadDir(yearPath)
		if err != nil {
			continue
		}

		for _, monthEntry := range monthEntries {
			if !monthEntry.IsDir() {
				continue
			}

			// Count files in month directory
			monthPath := filepath.Join(yearPath, monthEntry.Name())
			count := a.countFiles(monthPath)

			// Parse month name from directory
			monthName := a.extractMonthName(monthEntry.Name())
			if monthName == "" {
				continue
			}

			// Update month data
			if month, ok := yearData.Months[monthName]; ok {
				month.FileCount = count
				month.Status = getStatus(count)
				result.TotalFiles += count
			}
		}

		result.Years[year] = yearData
	}

	return result, nil
}

func (a *Analyzer) extractMonthName(dirName string) string {
	// Handle patterns like "01-January", "1-January", "January", etc.
	parts := strings.Split(dirName, "-")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// monthOrder returns the numerical order (1-12) for a given month name
// monthOrder returns the numerical order (1-12) for a given month name
func monthOrder(month string) int {
	monthOrders := map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}
	return monthOrders[month]
}

func (a *Analyzer) countFiles(path string) int {
	count := 0
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			count++
		}
		return nil
	})
	return count
}

func getStatus(count int) string {
    switch {
    case count == 0:
        return "missing"
    case count < 30:
        return "light"
    default:
        return "normal"
    }
}
