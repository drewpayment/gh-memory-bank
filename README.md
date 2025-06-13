# GitHub Memory Bank

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This repository contains a CLI tool to help users implement a memory bank system for use with AI agents like GitHub Copilot. The system is designed to maintain context across sessions using specialized modes that handle different phases of the development process.

## Overview

The Memory Bank system utilizes a structured approach to manage tasks, context, and progress throughout a development workflow. It relies on a set of predefined "chat modes," each tailored for a specific stage:

*   **VAN Mode (Initialization):** Sets up the initial context, checks the memory bank status, and determines the task's complexity level.
*   **PM Mode (Product Requirements):** Guides the creation of detailed Product Requirements Documents (PRDs) through structured clarifying questions and comprehensive documentation.
*   **PLAN Mode (Task Planning):** Creates a detailed plan for task execution based on the complexity.
*   **CREATIVE Mode (Design Decisions):** Facilitates detailed design and architecture work for components flagged during planning.
*   **IMPLEMENT Mode (Code Implementation):** Guides the building of planned changes, following the implementation plan and creative phase decisions.
*   **REFLECT+ARCHIVE Mode (Review & Archiving):** Facilitates reflection on the completed task and then archives relevant documentation, updating the Memory Bank.

## Key Features

*   **Mode-Specific Workflows:** Each mode has a defined process, often visualized with Mermaid diagrams, to guide the user and the AI.
*   **Contextual Memory:** The system uses a "Memory Bank" (a set of structured markdown files) to maintain context, track tasks (`tasks.md`), manage active context (`activeContext.md`), store creative decisions (`creative-*.md`), and log progress (`progress.md`).
*   **Isolation-Focused Design:** Each mode is designed to load only the rules and information it needs, optimizing for efficiency.
*   **Structured Documentation:** The system emphasizes clear documentation at each stage, including planning documents, design guidelines, implementation details, and reflection notes.
*   **Verification Checkpoints:** Each mode includes verification steps to ensure that processes are followed correctly and that the Memory Bank is updated appropriately.

## How It Works

The system is intended to be used with GitHub Copilot Chat's custom modes. The user interacts with the AI, and the AI follows the instructions and workflows defined in the `.chatmode.md` files.

1.  **Initialization (VAN Mode):** A new task begins in VAN mode, where the project brief and task complexity are established.
2.  **Product Requirements (PM Mode):** When needed, this mode creates comprehensive Product Requirements Documents (PRDs) through structured questioning, ensuring clear feature definition before planning begins.
3.  **Planning (PLAN Mode):** Based on the task's complexity, a detailed implementation plan is created. Components requiring significant design work are flagged for a "Creative Phase."
4.  **Creative Design (CREATIVE Mode):** If needed, this mode is used to explore design options (architecture, algorithms, UI/UX) for flagged components.
5.  **Implementation (IMPLEMENT Mode):** Code changes are made according to the plan and any creative design decisions.
6.  **Reflection and Archiving (REFLECT+ARCHIVE Mode):** After implementation, the work is reviewed, lessons are learned, and all relevant documentation is archived. The Memory Bank is updated to reflect the completed task.

## Installation

### For Users

Install via your preferred package manager:

```bash
# Homebrew (macOS/Linux)
brew install your-username/tap/gh-memory-bank

# Chocolatey (Windows)
choco install gh-memory-bank

# Direct download from GitHub Releases
# Download the appropriate binary for your platform from:
# https://github.com/your-username/gh-memory-bank/releases
```

### Manual Installation

```bash
# Install directly with Go
go install github.com/your-username/gh-memory-bank@latest
```

## Usage

Navigate to your project directory and run:

```bash
gh-memory-bank install
```

The tool will:
1. Ask you where to install the templates (defaults to current directory)
2. Create a `.github` directory with all the chatmode and instruction files
3. Install or merge a `.gitignore` file in your project root
4. Warn you if you're not in a git repository (but allow you to continue)

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Setting Up the Development Environment

```bash
# Clone the repository
git clone https://github.com/your-username/gh-memory-bank.git
cd gh-memory-bank

# Initialize Go module (if not already done)
go mod init github.com/your-username/gh-memory-bank

# Install dependencies (if any are added later)
go mod tidy
```

### Project Structure

```
gh-memory-bank/
├── go.mod                          # Go module definition
├── main.go                         # Main CLI application
├── main_test.go                    # Comprehensive test suite
├── README.md                       # This file
├── templates/                      # Template files to install
│   ├── chatmodes/                  # GitHub Copilot chat modes
│   │   ├── IMPLEMENT.chatmode.md
│   │   ├── CREATIVE.chatmode.md
│   │   ├── PLAN.chatmode.md
│   │   ├── PM.chatmode.md
│   │   ├── REFLECT.chatmode.md
│   │   └── VAN.chatmode.md
│   ├── instructions/               # Copilot instructions
│   │   └── main.instructions.md
│   └── gitignore                   # Gets renamed to .gitignore in target
└── .goreleaser.yml                 # GoReleaser config (for releases)
```

### Building and Testing

```bash
# Run tests
go test -v

# Run tests with coverage
go test -v -cover

# Run specific test
go test -v -run TestInstallFilesToPath

# Run benchmarks
go test -v -bench .

# Build for your platform
go build -o gh-memory-bank

# Test the built binary
./gh-memory-bank install
```

### Cross-Platform Building

```bash
# Build for Windows
GOOS=windows GOARCH=amd64 go build -o gh-memory-bank.exe

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o gh-memory-bank-mac

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o gh-memory-bank-linux
```

### Testing Your Changes

1. **Unit Tests**: Run `go test -v` to ensure all functionality works
2. **Manual Testing**: Build the binary and test it in a real project:
   ```bash
   go build -o gh-memory-bank
   cd /path/to/test/project
   /path/to/gh-memory-bank install
   ```
3. **Template Changes**: Modify files in `templates/` and test that they're properly embedded and installed

### Key Components

- **`main.go`**: Contains the CLI logic and file installation code
- **`templates/`**: All template files that get installed to user projects
- **Embedded Files**: Templates are embedded into the binary using Go's `//go:embed` feature
- **Gitignore Handling**: Special logic merges template gitignore with existing ones
- **Cross-Platform**: Uses `filepath` package for proper path handling across OS platforms

### Running Tests

The test suite covers:
- Basic file installation
- Git repository detection
- Directory creation
- Gitignore merging logic
- Content integrity verification
- Cross-platform path handling
- Performance benchmarks

```bash
# Run all tests
go test

# Verbose output
go test -v

# Run with race detection
go test -race

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Adding New Template Files

1. Add your new template file to the appropriate directory under `templates/`
2. The file will automatically be embedded and installed
3. Run tests to ensure everything works: `go test -v`
4. For special files (like gitignore), you may need to add custom handling logic in `main.go`

### Release Process

This project uses GoReleaser for automated releases:

```bash
# Tag a new version
git tag v1.0.0
git push origin v1.0.0

# GoReleaser will automatically:
# - Build binaries for multiple platforms
# - Create GitHub release
# - Generate Homebrew formula
# - Create Chocolatey package
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run the test suite: `go test -v`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Guidelines

- Write tests for new functionality
- Maintain cross-platform compatibility
- Update documentation for user-facing changes
- Follow Go best practices and formatting (`go fmt`)
- Ensure all tests pass before submitting

## License

[Add your license here]

---

This project provides the infrastructure to install and manage GitHub Copilot Memory Bank templates, making it easy for teams to adopt structured AI-assisted development workflows.