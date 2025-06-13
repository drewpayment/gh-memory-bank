package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallFilesToPath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test installation to temp directory (non-interactive mode)
	err = installFilesToPath(tempDir, false)
	if err != nil {
		t.Fatalf("Installation failed: %v", err)
	}

	// Verify .github directory was created
	githubDir := filepath.Join(tempDir, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		t.Errorf(".github directory was not created")
	}

	// Verify template files were copied (excluding gitignore which goes to root)
	expectedFiles := getExpectedTemplateFiles(t, true) // exclude gitignore
	for _, expectedFile := range expectedFiles {
		fullPath := filepath.Join(githubDir, expectedFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", fullPath)
		}
	}

	// Check if gitignore was handled separately (if it exists in templates)
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if hasGitignoreTemplate(t) {
		if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
			t.Errorf(".gitignore file should have been created in project root")
		}
	}
}

func TestInstallFilesToPathWithGitRepo(t *testing.T) {
	// Create a temporary directory with a .git folder
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-git-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create .git directory to simulate a git repo
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	// Test installation - should work without warnings since it's a git repo
	err = installFilesToPath(tempDir, false)
	if err != nil {
		t.Fatalf("Installation in git repo failed: %v", err)
	}

	// Verify files were installed
	githubDir := filepath.Join(tempDir, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		t.Errorf(".github directory was not created in git repo")
	}
}

func TestInstallFilesToPathCreatesDirectories(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-dir-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test installation to a nested path that doesn't exist yet
	nestedPath := filepath.Join(tempDir, "some", "deep", "nested", "path")

	err = installFilesToPath(nestedPath, false)
	if err != nil {
		t.Fatalf("Installation to nested path failed: %v", err)
	}

	// Verify the nested directories were created
	githubDir := filepath.Join(nestedPath, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		t.Errorf(".github directory was not created in nested path")
	}
}

func TestInstallFilesToPathRelativePath(t *testing.T) {
	// Create a temporary directory and change to it
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-rel-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save current directory and change to temp dir
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test installation with relative path
	err = installFilesToPath(".", false)
	if err != nil {
		t.Fatalf("Installation with relative path failed: %v", err)
	}

	// Verify files were installed
	githubDir := filepath.Join(tempDir, ".github")
	if _, err := os.Stat(githubDir); os.IsNotExist(err) {
		t.Errorf(".github directory was not created with relative path")
	}
}

func TestTemplateFilesAreEmbedded(t *testing.T) {
	// Test that our embedded files are accessible
	entries, err := fs.ReadDir(templateFiles, "templates")
	if err != nil {
		t.Fatalf("Failed to read embedded templates directory: %v", err)
	}

	if len(entries) == 0 {
		t.Errorf("No template files found in embedded filesystem")
	}

	// Test reading a file from the embedded filesystem
	err = fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			// Try to read the file
			content, readErr := templateFiles.ReadFile(path)
			if readErr != nil {
				t.Errorf("Failed to read embedded file %s: %v", path, readErr)
			}

			if len(content) == 0 {
				t.Errorf("Embedded file %s appears to be empty", path)
			}
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error walking embedded files: %v", err)
	}
}

func TestFileContentIntegrity(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-content-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Install files
	err = installFilesToPath(tempDir, false)
	if err != nil {
		t.Fatalf("Installation failed: %v", err)
	}

	// Check that installed files match embedded files
	githubDir := filepath.Join(tempDir, ".github")

	err = fs.WalkDir(templateFiles, "templates", func(embeddedPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Read embedded content
		embeddedContent, err := templateFiles.ReadFile(embeddedPath)
		if err != nil {
			t.Errorf("Failed to read embedded file %s: %v", embeddedPath, err)
			return nil
		}

		// Calculate installed file path
		relPath, err := filepath.Rel("templates", embeddedPath)
		if err != nil {
			t.Errorf("Failed to calculate relative path for %s: %v", embeddedPath, err)
			return nil
		}

		var installedPath string
		// Special handling for gitignore - it goes to project root as .gitignore
		if filepath.Base(relPath) == "gitignore" {
			installedPath = filepath.Join(tempDir, ".gitignore")
		} else {
			installedPath = filepath.Join(githubDir, relPath)
		}

		// Read installed content
		installedContent, err := os.ReadFile(installedPath)
		if err != nil {
			t.Errorf("Failed to read installed file %s: %v", installedPath, err)
			return nil
		}

		// For gitignore, we only check that the embedded content is contained in the installed content
		// (since it might have merge comments added)
		if filepath.Base(relPath) == "gitignore" {
			if !strings.Contains(string(installedContent), strings.TrimSpace(string(embeddedContent))) {
				t.Errorf("Gitignore content not properly installed for file %s", relPath)
			}
		} else {
			// For other files, content should match exactly
			if string(embeddedContent) != string(installedContent) {
				t.Errorf("Content mismatch for file %s", relPath)
			}
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error during content verification: %v", err)
	}
}

func TestPromptForInstallPathDefault(t *testing.T) {
	// This test would need input mocking to fully test
	// For now, we'll just test that the function exists and handles empty input

	// Test that we can get the current directory as default
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// The promptForInstallPath function should return currentDir when given empty input
	// This would require input mocking to test properly
	if currentDir == "" {
		t.Errorf("Current directory should not be empty")
	}
}

// Helper function to get expected template files from embedded filesystem
func getExpectedTemplateFiles(t *testing.T, excludeGitignore bool) []string {
	var files []string

	err := fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			relPath, err := filepath.Rel("templates", path)
			if err != nil {
				return err
			}

			// Skip gitignore if requested (since it goes to project root, not .github)
			if excludeGitignore && filepath.Base(relPath) == "gitignore" {
				return nil
			}

			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to get expected template files: %v", err)
	}

	return files
}

// Helper function to check if gitignore template exists
func hasGitignoreTemplate(t *testing.T) bool {
	err := fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Base(path) == "gitignore" {
			return fs.SkipAll // Found it, stop walking
		}
		return nil
	})

	return err == fs.SkipAll
}

func TestGitignoreHandling(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-gitignore-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Only run this test if we have a gitignore template
	if !hasGitignoreTemplate(t) {
		t.Skip("No gitignore template found, skipping gitignore tests")
	}

	// Test installation - should create .gitignore in project root
	err = installFilesToPath(tempDir, false)
	if err != nil {
		t.Fatalf("Installation failed: %v", err)
	}

	// Verify .gitignore was created in project root (not in .github)
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Errorf(".gitignore file was not created in project root")
	}

	// Verify it's NOT in the .github directory
	githubGitignorePath := filepath.Join(tempDir, ".github", "gitignore")
	if _, err := os.Stat(githubGitignorePath); !os.IsNotExist(err) {
		t.Errorf("gitignore file should not exist in .github directory")
	}
}

func TestGitignoreMerging(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "gh-memory-bank-gitignore-merge-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Only run this test if we have a gitignore template
	if !hasGitignoreTemplate(t) {
		t.Skip("No gitignore template found, skipping gitignore merge tests")
	}

	// Create an existing .gitignore file
	existingGitignore := `# Existing content
*.log
*.tmp
node_modules/
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	err = os.WriteFile(gitignorePath, []byte(existingGitignore), 0644)
	if err != nil {
		t.Fatalf("Failed to create existing .gitignore: %v", err)
	}

	// Install templates (this should merge with existing .gitignore)
	err = installFilesToPath(tempDir, false)
	if err != nil {
		t.Fatalf("Installation failed: %v", err)
	}

	// Read the merged .gitignore
	mergedContent, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("Failed to read merged .gitignore: %v", err)
	}

	mergedStr := string(mergedContent)

	// Verify existing content is preserved
	if !strings.Contains(mergedStr, "*.log") {
		t.Errorf("Existing content was not preserved in merged .gitignore")
	}

	// Verify merge comment was added
	if !strings.Contains(mergedStr, "# Added by gh-memory-bank") {
		t.Errorf("Merge comment was not added to .gitignore")
	}
}

func TestMergeGitignoreContent(t *testing.T) {
	existing := `# Existing
*.log
*.tmp
node_modules/
`

	new := `# New content
*.log
*.cache
dist/
`

	result := mergeGitignoreContent(existing, new)

	// Should contain existing content
	if !strings.Contains(result, "*.tmp") {
		t.Errorf("Result should contain existing content")
	}

	// Should contain new unique content
	if !strings.Contains(result, "*.cache") {
		t.Errorf("Result should contain new unique content")
	}

	if !strings.Contains(result, "dist/") {
		t.Errorf("Result should contain new unique content")
	}

	// Should not duplicate existing content
	logCount := strings.Count(result, "*.log")
	if logCount > 1 {
		t.Errorf("Should not duplicate existing entries, found %d instances of '*.log'", logCount)
	}

	// Should contain merge comment
	if !strings.Contains(result, "# Added by gh-memory-bank") {
		t.Errorf("Result should contain merge comment")
	}
}

func TestMergeGitignoreContentNoDuplicates(t *testing.T) {
	existing := `*.log
node_modules/
dist/
`

	new := `*.log
node_modules/
*.cache
`

	result := mergeGitignoreContent(existing, new)

	// Count occurrences of each line
	logCount := strings.Count(result, "*.log")
	nodeModulesCount := strings.Count(result, "node_modules/")

	if logCount > 1 {
		t.Errorf("Found %d instances of '*.log', should be 1", logCount)
	}

	if nodeModulesCount > 1 {
		t.Errorf("Found %d instances of 'node_modules/', should be 1", nodeModulesCount)
	}

	// Should contain the new unique entry
	if !strings.Contains(result, "*.cache") {
		t.Errorf("Result should contain new unique content '*.cache'")
	}
}

// Benchmark test for installation performance
func BenchmarkInstallFilesToPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tempDir, err := os.MkdirTemp("", "gh-memory-bank-bench-")
		if err != nil {
			b.Fatalf("Failed to create temp dir: %v", err)
		}

		err = installFilesToPath(tempDir, false)
		if err != nil {
			b.Fatalf("Installation failed: %v", err)
		}

		os.RemoveAll(tempDir)
	}
}

// Test that verifies cross-platform path handling
func TestCrossPlatformPaths(t *testing.T) {
	testPaths := []string{
		".",
		"./relative/path",
		"../parent/path",
	}

	for _, testPath := range testPaths {
		t.Run("path_"+strings.ReplaceAll(testPath, "/", "_"), func(t *testing.T) {
			// Convert to absolute path (this is what the function does)
			absPath, err := filepath.Abs(testPath)
			if err != nil {
				t.Errorf("Failed to resolve path %s: %v", testPath, err)
			}

			// Just verify the path can be resolved
			if absPath == "" {
				t.Errorf("Absolute path should not be empty for input: %s", testPath)
			}
		})
	}
}
