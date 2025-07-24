package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"terraform-provider-yamlflattener/internal/utils"
)

func TestParseYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
	}{
		{
			name:        "Valid simple YAML",
			yamlContent: "key: value",
			wantErr:     false,
		},
		{
			name: "Valid complex YAML",
			yamlContent: `
key1: value1
key2:
  nested: value2
array:
  - item1
  - item2
`,
			wantErr: false,
		},
		{
			name:        "Empty YAML",
			yamlContent: "",
			wantErr:     true,
		},
		{
			name:        "Invalid YAML",
			yamlContent: "key: : value",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ParseYAML(tt.yamlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Errorf("utils.ParseYAML() returned nil result for valid YAML")
			}
		})
	}
}

func TestValidateYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
	}{
		{
			name:        "Valid YAML",
			yamlContent: "key: value",
			wantErr:     false,
		},
		{
			name:        "Empty YAML",
			yamlContent: "",
			wantErr:     true,
		},
		{
			name:        "Invalid YAML syntax",
			yamlContent: "key: : value",
			wantErr:     true,
		},
		{
			name:        "Invalid YAML structure",
			yamlContent: "key:\n  - item1\n  item2",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateYAML(tt.yamlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.ValidateYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadYAMLFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create a valid YAML file
	validFilePath := filepath.Join(tempDir, "valid.yaml")
	validContent := "key: value\narray:\n  - item1\n  - item2"
	err = os.WriteFile(validFilePath, []byte(validContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a directory with the same name as a potential YAML file
	dirPath := filepath.Join(tempDir, "dir.yaml")
	err = os.Mkdir(dirPath, 0750)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid file",
			filePath: validFilePath,
			want:     validContent,
			wantErr:  false,
		},
		{
			name:     "Non-existent file",
			filePath: filepath.Join(tempDir, "nonexistent.yaml"),
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Directory instead of file",
			filePath: dirPath,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Empty file path",
			filePath: "",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Path with directory traversal",
			filePath: "../../../etc/passwd",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ReadYAMLFile(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.ReadYAMLFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("utils.ReadYAMLFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAndParseYAMLFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create a valid YAML file
	validFilePath := filepath.Join(tempDir, "valid.yaml")
	validContent := "key: value\narray:\n  - item1\n  - item2"
	err = os.WriteFile(validFilePath, []byte(validContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create an invalid YAML file
	invalidFilePath := filepath.Join(tempDir, "invalid.yaml")
	invalidContent := "key: : value"
	err = os.WriteFile(invalidFilePath, []byte(invalidContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		wantErr  bool
	}{
		{
			name:     "Valid YAML file",
			filePath: validFilePath,
			wantErr:  false,
		},
		{
			name:     "Invalid YAML file",
			filePath: invalidFilePath,
			wantErr:  true,
		},
		{
			name:     "Non-existent file",
			filePath: filepath.Join(tempDir, "nonexistent.yaml"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.ReadAndParseYAMLFile(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("utils.ReadAndParseYAMLFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Errorf("utils.ReadAndParseYAMLFile() returned nil result for valid YAML file")
			}
		})
	}
}
