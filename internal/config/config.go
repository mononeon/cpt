package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	DefaultLanguage string            `json:"default_language"`
	NamingFormat    string            `json:"naming_format"` // veryshortform, shortform, longform
	CompileCommands map[string]string `json:"compile_commands"`
	RunCommands     map[string]string `json:"run_commands"`
	Templates       map[string]string `json:"templates"`
}

var GlobalConfig Config

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cptconfig"), nil
}

// StripJSONComments removes single-line (//) and multi-line (/* */) comments from a JSONC string.
func StripJSONComments(data []byte) []byte {
	// Remove /* ... */
	reMultiline := regexp.MustCompile(`(?s)/\*.*?\*/`)
	data = reMultiline.ReplaceAll(data, nil)
	// Remove // ...
	reSingleline := regexp.MustCompile(`//.*`)
	data = reSingleline.ReplaceAll(data, nil)
	return data
}

func LoadConfig() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Set defaults
	GlobalConfig = Config{
		DefaultLanguage: "cpp",
		NamingFormat:    "veryshortform",
		CompileCommands: map[string]string{
			"cpp": "g++ -O2 -Wall -Wextra -std=c++17 {source} -o {output}",
		},
		RunCommands: map[string]string{
			"cpp": "{executable}",
			"py":  "python {source}",
		},
		Templates: map[string]string{
			"cpp": "",
			"py":  "",
		},
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Create empty config if it doesn't exist
		emptyCfg := []byte("// CPT Configuration File\n// You can override default settings here.\n{\n}\n")
		return os.WriteFile(path, emptyCfg, 0644)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	data = StripJSONComments(data)

	// Merge with defaults
	var userConfig Config
	if err := json.Unmarshal(data, &userConfig); err != nil {
		return fmt.Errorf("error parsing .cptconfig: %v", err)
	}

	if userConfig.DefaultLanguage != "" {
		GlobalConfig.DefaultLanguage = userConfig.DefaultLanguage
	}
	if userConfig.NamingFormat != "" {
		GlobalConfig.NamingFormat = userConfig.NamingFormat
	}
	for k, v := range userConfig.CompileCommands {
		GlobalConfig.CompileCommands[k] = v
	}
	for k, v := range userConfig.RunCommands {
		GlobalConfig.RunCommands[k] = v
	}
	for k, v := range userConfig.Templates {
		GlobalConfig.Templates[k] = v
	}

	return nil
}
