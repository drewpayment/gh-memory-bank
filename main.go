package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Embed all files from the templates directory
//
//go:embed templates/*
var templateFiles embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gh-memory-bank install")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		if err := installFiles(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully installed GitHub Copilot memory bank templates!")
	default:
		fmt.Println("Unknown command. Use: gh-memory-bank install")
	}
}

func installFiles() error {
	// Ask user for installation path
	installPath, err := promptForInstallPath()
	if err != nil {
		return fmt.Errorf("failed to get installation path: %w", err)
	}

	return installFilesToPath(installPath, true)
}

func installFilesToPath(installPath string, interactive bool) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(installPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	fmt.Printf("Installing to: %s\n", absPath)

	// Check if the target directory exists or create it
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Check if we're in a git repo (relative to the target path)
	gitPath := filepath.Join(absPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		if interactive {
			fmt.Printf("Warning: %s doesn't appear to be a git repository\n", absPath)
			if !promptForConfirmation("Continue anyway?") {
				return fmt.Errorf("installation cancelled")
			}
		}
	}

	// Create .github directory in the target path
	githubDir := filepath.Join(absPath, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		return fmt.Errorf("failed to create .github directory: %w", err)
	}

	// Walk through embedded files and copy them
	return fs.WalkDir(templateFiles, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Read the embedded file
		content, err := templateFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Calculate destination path (remove "templates/" prefix)
		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		// Special handling for gitignore file
		if filepath.Base(relPath) == "gitignore" {
			return handleGitignoreFile(absPath, content)
		}

		// Regular file handling - goes to .github directory
		destPath := filepath.Join(githubDir, relPath)

		// Create destination directory if needed
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		// Write the file
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		fmt.Printf("Installed: %s\n", destPath)
		return nil
	})
}

func promptForInstallPath() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	fmt.Printf("Enter installation path (default: %s): ", currentDir)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return currentDir, nil
	}

	return input, nil
}

func promptForConfirmation(message string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s (y/N): ", message)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func handleGitignoreFile(projectRoot string, newContent []byte) error {
	gitignorePath := filepath.Join(projectRoot, ".gitignore")

	// Check if .gitignore already exists
	existingContent, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing .gitignore: %w", err)
	}

	var finalContent string

	if os.IsNotExist(err) {
		// No existing .gitignore, just use the new content
		finalContent = string(newContent)
		fmt.Printf("Created: %s\n", gitignorePath)
	} else {
		// Merge existing and new content
		finalContent = mergeGitignoreContent(string(existingContent), string(newContent))
		fmt.Printf("Updated: %s (merged with existing content)\n", gitignorePath)
	}

	// Write the final content
	if err := os.WriteFile(gitignorePath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	return nil
}

func mergeGitignoreContent(existing, new string) string {
	// Parse both files into lines
	existingLines := strings.Split(strings.TrimSpace(existing), "\n")
	newLines := strings.Split(strings.TrimSpace(new), "\n")

	// Create a set of existing lines to avoid duplicates
	existingSet := make(map[string]bool)
	for _, line := range existingLines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			existingSet[trimmed] = true
		}
	}

	// Start with existing content
	result := existing

	// Add new lines that don't already exist
	var newUniqueLines []string
	for _, line := range newLines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !existingSet[trimmed] {
			newUniqueLines = append(newUniqueLines, line)
		}
	}

	// If we have new unique lines, add them
	if len(newUniqueLines) > 0 {
		// Ensure existing content ends with a newline
		if !strings.HasSuffix(result, "\n") {
			result += "\n"
		}

		// Add a comment to indicate where the new content starts
		result += "\n# Added by gh-memory-bank\n"
		result += strings.Join(newUniqueLines, "\n")
		result += "\n"
	}

	return result
}
