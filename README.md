# Git Fame Analyzer

A Go-based command-line utility for detailed analysis of author contributions to a git repository. The program calculates the number of lines, commits, and affected files for each participant at a specific point in the project's history.

## Description

`git-fame` is a command-line tool that helps visualize the "fame" of developers. The utility scans repository files, matches each line of code to the last commit that changed it via `git blame`, and aggregates data by authors or committers.

**Features:**

* **Flexible filtering**: Restrict analysis by file extensions (`.go`, `.md`), programming languages, or specific paths.
* **Exclusion support**: Use Glob patterns to ignore folders (e.g., `vendor/*`) or specific files.
* **Variety of formats**: Output results as a readable table (tabular), CSV, JSON, or JSON-lines for subsequent automated processing.
* **High accuracy**: Accounts for the repository state at any specified revision (SHA-1, branch, tag).
* **Sorting**: Ability to rank authors by the number of lines, commits, or files.

## Usage

### Basic Commands

```bash
# Basic run in the current repository (default tabular output)
gitfame

# Analyze only Go and Markdown files, sort by number of commits
gitfame --extensions='.go,.md' --order-by=commits

# Output in JSON format for integration with other utilities
gitfame --format=json --exclude='vendor/**,**/*_test.go'
```

### Flag Reference

| Flag | Default Value | Description |
| --- | --- | --- |
| `--repository` | `.` | Path to the Git repository. |
| `--revision` | `HEAD` | Pointer to a commit (branch, hash, or tag). |
| `--order-by` | `lines` | Sorting key: `lines`, `commits`, or `files`. |
| `--format` | `tabular` | Output format: `tabular`, `csv`, `json`, `json-lines`. |
| `--extensions` |  | List of extensions separated by commas (e.g., `.go,.js`). |
| `--languages` |  | List of languages (e.g., `go,python`). |
| `--use-committer` | `false` | Use committer instead of author. |

## Installation

```bash
git clone https://github.com/mhgffqwoer/git-fame.git
cd git-fame
go install ./cmd/gitfame

```

## How It Works

The analysis process is divided into several stages:

1. **File Collection**: Using `git ls-tree`, a list of files relevant to the specified revision is formed.
2. **Filtering**: User filters for extensions, languages, and Glob patterns are applied.
3. **Blame Analysis**: For each file that passes the filter, `git blame --porcelain` is called. This format effectively groups adjacent lines from the same commit.
4. **Aggregation**: The program matches commit data with authors (or committers when the `--use-committer` flag is used) and accumulates three metrics:
* **Lines**: Total number of lines belonging to the author.
* **Commits**: Number of unique commits by the author in the selected subset of files.
* **Files**: Number of unique files to which the author contributed changes.

## Project Structure

```text
/ (root)
├─ cmd/
│  └─ gitfame/              # Main entry point
│     └─ main.go            # CLI argument parsing and analyzer initialization
├─ configs/
│  ├─ embed.go              # Mechanism for embedding static files
│  └─ language_extensions.json # Mapping of extensions to languages
├─ internal/
│  ├─ analyzer/             # Analysis system core
│  │  ├─ entities/          # Data models (files, history)
│  │  ├─ filter/            # File filtering logic (Glob, extensions)
│  │  ├─ git_adapters/      # Adapters for streaming reads from Git (blame, ls-tree)
│  │  └─ storage/           # Temporary storage for intermediate results
│  ├─ git_client/           # Low-level wrapper over Git commands
│  ├─ metrics/              # Structures for counting lines, commits, and files
│  ├─ output/               # Result generation module
│  │  ├─ formatter/         # Formatting: CSV, JSON, Tabular
│  │  └─ sorter/            # Sorting logic by keys (lines, commits, files)
│  └─ shared/               # General utility functions
└─ test/
   └─ integration/          # Integration tests on real git bundles
```