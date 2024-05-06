package file

import (
	"os"
	"strings"
	"testing"
)

// Test cases for successful scenarios
func TestGetProjectDir_Success(t *testing.T) {
	// Test case 1: Project name exists in the current working directory
	// Replace 'projectName' with an existing project name on your system.
	projectName := "currency-service"
	actualPath, err := GetProjectDir(projectName)
	if err != nil {
		t.Errorf("Expected success, but got an error: %v", err)
	}
	if !strings.Contains(actualPath, projectName) {
		t.Errorf("Expected containes: %s, but got: %s", projectName, actualPath)
	}

	// Test case 2: Empty project name (should return an error)
	emptyProjectName := "not_found"
	_, err = GetProjectDir(emptyProjectName)
	if err == nil {
		t.Error("Expected an error for an empty project name, but got nil")
	}
}

// Test cases for error scenarios
func TestGetProjectDir_Errors(t *testing.T) {
	// Test case 1: Project directory not found in the current working directory
	projectName := "non_existent_project"
	_, err := GetProjectDir(projectName)
	if err != ErrProjectDirectoryNotFound {
		t.Errorf("Expected error 'ErrProjectDirectoryNotFound', but got: %v", err)
	}

	// Test case 2: Getwd function error
	// Temporarily change the current working directory to a non-existent path to simulate a Getwd error.
	previousDir, _ := os.Getwd()
	defer os.Chdir(previousDir) // Restore the previous directory.
	os.Chdir("non_existent_path")
	_, err = GetProjectDir(projectName)
	if err == nil {
		t.Error("Expected an error due to Getwd failure, but got nil")
	}
}

func TestFirstIndexOfString(t *testing.T) {
	tests := []struct {
		source   string
		value    string
		expected int
	}{
		{"abcdefgh", "def", 3},
		{"abcabcabc", "abc", 0},
		{"xyz", "abc", -1}, // Not found case
		{"", "test", -1},   // Empty source case
		{"test", "", 0},    // Empty value case
		{"aaa", "aa", 0},   // Overlapping case
	}

	for _, test := range tests {
		result := firstIndexOfString(test.source, test.value)
		if result != test.expected {
			t.Errorf("For source=%s, value=%s, expected %d but got %d", test.source, test.value, test.expected, result)
		}
	}
}
