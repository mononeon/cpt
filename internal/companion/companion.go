package companion

import (
	"cpt/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Payload struct {
	Name  string     `json:"name"`
	Group string     `json:"group"`
	URL   string     `json:"url"`
	Tests []TestCase `json:"tests"`
}

func sanitizeName(name string) string {
	// Remove illegal path characters
	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = re.ReplaceAllString(name, "")
	return strings.TrimSpace(name)
}

func parseProblemName(payload Payload, format string) string {
	// Fallback to simple name if format is unknown or parsing fails
	if format == "" {
		format = config.GlobalConfig.NamingFormat
	}

	if format == "longform" {
		return sanitizeName(payload.Name)
	}

	// Try to extract contest and index from URL
	// Codeforces example: https://codeforces.com/contest/1337/problem/A
	// or https://codeforces.com/problemset/problem/1337/A
	contest := ""
	index := ""
	
	reCF1 := regexp.MustCompile(`contest/(\d+)/problem/([A-Z0-9]+)`)
	reCF2 := regexp.MustCompile(`problem/(\d+)/([A-Z0-9]+)`)

	if m := reCF1.FindStringSubmatch(payload.URL); len(m) == 3 {
		contest = m[1]
		index = m[2]
	} else if m := reCF2.FindStringSubmatch(payload.URL); len(m) == 3 {
		contest = m[1]
		index = m[2]
	}

	if index == "" {
		// Fallback: take first word of the name
		parts := strings.Fields(payload.Name)
		if len(parts) > 0 {
			index = sanitizeName(parts[0])
		} else {
			index = "Unknown"
		}
	}

	// Remove trailing dots (e.g. if name is "A. Problem")
	index = strings.TrimSuffix(index, ".")

	if format == "veryshortform" {
		return index
	}

	if format == "shortform" {
		if contest != "" {
			return contest + index
		}
		return index // fallback
	}

	return index // fallback for unknown format
}

func getTemplateContent(ext string) string {
	templatePath := config.GlobalConfig.Templates[ext]
	if templatePath != "" {
		data, err := os.ReadFile(templatePath)
		if err == nil {
			return string(data)
		}
	}
	return ""
}

func createProblemStructure(problemName string, tests []TestCase) {
	ext := config.GlobalConfig.DefaultLanguage
	sourceFile := fmt.Sprintf("%s.%s", problemName, ext)

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		content := getTemplateContent(ext)
		os.WriteFile(sourceFile, []byte(content), 0644)
		fmt.Printf("\033[32mCreated %s\033[0m\n", sourceFile)
	} else {
		fmt.Printf("\033[33m%s already exists, skipping...\033[0m\n", sourceFile)
	}

	cptDir := filepath.Join(".cpt", problemName)
	os.MkdirAll(cptDir, 0755)

	for i, test := range tests {
		inFile := filepath.Join(cptDir, fmt.Sprintf("in_%d.txt", i+1))
		outFile := filepath.Join(cptDir, fmt.Sprintf("out_%d.txt", i+1))

		os.WriteFile(inFile, []byte(test.Input), 0644)
		os.WriteFile(outFile, []byte(test.Output), 0644)
	}

	fmt.Printf("\033[32mSaved %d test cases for %s\033[0m\n", len(tests), problemName)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	var payload Payload
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Printf("\033[31mError parsing payload: %v\033[0m\n", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	problemName := parseProblemName(payload, "")
	fmt.Printf("\033[1;34mReceived problem:\033[0m %s\n", problemName)
	createProblemStructure(problemName, payload.Tests)

	w.WriteHeader(http.StatusOK)
}

func Listen(port int, format string) {
	if format != "" {
		config.GlobalConfig.NamingFormat = format
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	http.HandleFunc("/", handler)
	fmt.Printf("\033[32mListening for Competitive Companion on %s...\033[0m\n", addr)
	fmt.Println("Press Ctrl+C to stop.")

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("\033[31mCould not start server on %s: %v\033[0m\n", addr, err)
	}
}
